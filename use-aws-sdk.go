package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/eks"

	"fmt"
	"os"
)

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1")},
	)

	clientEks := eks.New(sess)

	eksDescription, err := clientEks.DescribeCluster(&eks.DescribeClusterInput{
		Name: aws.String("bla"),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case eks.ErrCodeResourceNotFoundException:
				fmt.Println(eks.ErrCodeResourceNotFoundException, aerr.Error())
			case eks.ErrCodeClientException:
				fmt.Println(eks.ErrCodeClientException, aerr.Error())
			case eks.ErrCodeServerException:
				fmt.Println(eks.ErrCodeServerException, aerr.Error())
			case eks.ErrCodeServiceUnavailableException:
				fmt.Println(eks.ErrCodeServiceUnavailableException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}

	// for _, subnets := range eksDescription.Cluster.ResourcesVpcConfig.SubnetIds {
	// 	fmt.Println(*subnets)
	// }

	clientEc2 := ec2.New(sess)
	inputSubnets := &ec2.DescribeSubnetsInput{
		SubnetIds: eksDescription.Cluster.ResourcesVpcConfig.SubnetIds,
	}

	subnetsDescription, err := clientEc2.DescribeSubnets(inputSubnets)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}

	// fmt.Println(subnetsDescription)
	for _, subnets := range subnetsDescription.Subnets {
		fmt.Println(*subnets.AvailableIpAddressCount, *subnets.CidrBlock)
	}
}
