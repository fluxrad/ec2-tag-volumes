package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var InstanceID = flag.String("i", "", "The instance id whose volumes to tag")
var DryRun = flag.Bool("d", false, "Perform a dry run. Don't make any changes")
var Region = flag.String("r", "us-east-1", "The AWS region to use")
var AWSConfig *aws.Config

func init() {
	AWSConfig = &aws.Config{
		Region: aws.String(*Region),
	}
}

// Validate ID validates that the a given ec2 instance id is valid
func ValidateID(i *string) (bool, error) {
	m, err := regexp.MatchString("^i-[a-z0-9]{8}$", *i)
	if err != nil {
		return false, err
	}

	return m, nil
}

// DescribeInstance gets an Instance object from AWS and returns it, stripping
// out Reservation data.
func DescribeInstance(c *ec2.EC2, i *string) (*ec2.Instance, error) {
	q := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{i},
	}

	resp, err := c.DescribeInstances(q)
	if err != nil {
		return nil, err
	}

	if len(resp.Reservations[0].Instances) != 1 {
		return nil, errors.New("Found more than one instance. Bailing")
	}

	return resp.Reservations[0].Instances[0], nil
}

// NameTag returns the Name tag for an instance
func NameTag(i *ec2.Instance) (*string, error) {
	var n *string

	for _, t := range i.Tags {
		if *t.Key == "Name" {
			n = t.Value
			return n, nil
		}
	}

	// We didn't find a name tag
	return nil, errors.New("Could not find Name tag")
}

// TagVolumesForInstance tags the volumes attached to a given instance with a
// specified string and the device name
func TagVolumesForInstance(c *ec2.EC2, i *ec2.Instance, n *string) error {
	for _, m := range i.BlockDeviceMappings {
		dn := m.DeviceName
		e := m.Ebs.VolumeId
		fmt.Printf("Tag volume %s:  %s - %s\n", *e, *n, *dn)

		p := &ec2.CreateTagsInput{
			Resources: []*string{
				aws.String(*e),
			},
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String(fmt.Sprintf("%s - %s", *n, *dn)),
				},
			},
			DryRun: aws.Bool(*DryRun),
		}

		resp, err := c.CreateTags(p)
		if err != nil {
			// Don't return because apparently when `DryRun` is set to true,
			// the API object also comes with an error.
			fmt.Println(err)
		}
		fmt.Println(resp)
	}

	return nil
}

func main() {
	flag.Parse()

	ok, err := ValidateID(InstanceID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !ok {
		fmt.Println("InstanceID doesn't appear to be valid")
		os.Exit(1)
	}

	conn := ec2.New(AWSConfig)

	instance, err := DescribeInstance(conn, InstanceID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	name, err := NameTag(instance)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Found Name tag: %s\n", *name)

	err = TagVolumesForInstance(conn, instance, name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
