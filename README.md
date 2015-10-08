ec2-tag-volumes
===============

A simple app to tag volumes attached to an EC2 insatnce. By default, volume are tagged with the instance Name tag, and device name

### Usage

```
./ec2-tag-volumes -i i-12345678
```


Given an ec2 instance with Name: `foo` with volumes `xvda` and `xvdf`, ec2-tag-volumes will tag them:

* `foo - /dev/xvda`
* `foo - /dev/xvdf`

