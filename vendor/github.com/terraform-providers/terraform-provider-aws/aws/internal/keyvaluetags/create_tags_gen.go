// Code generated by generators/createtags/main.go; DO NOT EDIT.

package keyvaluetags

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const EventualConsistencyTimeout = 5 * time.Minute

// Similar to isAWSErr from aws/awserr.go
// TODO: Add and export in shared package
func isAWSErrCode(err error, code string) bool {
	var awsErr awserr.Error
	if errors.As(err, &awsErr) {
		return awsErr.Code() == code
	}
	return false
}

// TODO: Add and export in shared package
func isAWSErrCodeContains(err error, code string) bool {
	var awsErr awserr.Error
	if errors.As(err, &awsErr) {
		return strings.Contains(awsErr.Code(), code)
	}
	return false
}

// Copied from aws/utils.go
// TODO: Export in shared package or add to Terraform Plugin SDK
func isResourceTimeoutError(err error) bool {
	timeoutErr, ok := err.(*resource.TimeoutError)
	return ok && timeoutErr.LastError == nil
}

// Ec2CreateTags creates ec2 service tags for new resources.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func Ec2CreateTags(conn *ec2.EC2, identifier string, tagsMap interface{}) error {
	tags := New(tagsMap)
	input := &ec2.CreateTagsInput{
		Resources: aws.StringSlice([]string{identifier}),
		Tags:      tags.IgnoreAws().Ec2Tags(),
	}

	err := resource.Retry(EventualConsistencyTimeout, func() *resource.RetryError {
		_, err := conn.CreateTags(input)

		if isAWSErrCodeContains(err, ".NotFound") {
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if isResourceTimeoutError(err) {
		_, err = conn.CreateTags(input)
	}

	if err != nil {
		return fmt.Errorf("error tagging resource (%s): %w", identifier, err)
	}

	return nil
}
