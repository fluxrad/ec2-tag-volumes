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

var InstanceID = flag.String("i", "", "the instance id whose volumes to tag")
var AWSConfig *aws.Config

func init() {
	AWSConfig = &aws.Config{
		Region: aws.String("us-east-1"),
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

func instanceName(c *ec2.EC2, i *string) (*string, error) {
	q := &ec2.DescribeTagsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("key"),
				Values: []*string{
					aws.String("Name"),
				},
			},
			{
				Name:   aws.String("resource-id"),
				Values: []*string{i},
			},
		},
	}

	resp, err := c.DescribeTags(q)
	if err != nil {
		return nil, err
	}

	if len(resp.Tags) != 1 {
		return nil, errors.New("Found more than one name tag.")
	}

	return resp.Tags[0].Value, nil
}

func volumesForInstance(c *ec2.EC2, i *string) ([]*string, error) {
	// blah!

	return nil, nil
}

func main() {
	flag.Parse()

	valid, err := ValidateID(InstanceID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if valid == false {
		fmt.Println("Instance ID doesn't appear to be valid")
		os.Exit(1)
	}

	conn := ec2.New(AWSConfig)

	name, err := instanceName(conn, InstanceID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Found name for instance %s: %s\n", *InstanceID, *name)

	vols, err := volumesForInstance(conn, InstanceID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)

	}
}
