package main

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
)

type eksService struct {
	sess  *session.Session
	svc   *eks.EKS
	param *eksParameters
}

type eksParameters struct {
	clustername string
	endpoint    string
	cadata      string
}

func NewEksService(region string) *eksService {
	sess := session.Must(session.NewSession())
	svc := eks.New(sess, aws.NewConfig().WithRegion(region))
	es := &eksService{
		sess:  sess,
		svc:   svc,
		param: new(eksParameters),
	}
	return es
}

func (es *eksService) ReadParameters(clustername string) error {
	// Extract the target cluster information by name.
	input := &eks.DescribeClusterInput{
		Name: aws.String(clustername),
	}
	result, err := es.svc.DescribeCluster(input)
	if err != nil {
		return errors.New(fmt.Sprintf("cannot access to AWS resources\n(detail) %s", err.Error()))
	}

	// Read and set the required properties.
	es.param.clustername = clustername
	es.param.endpoint = *result.Cluster.Endpoint
	es.param.cadata = *result.Cluster.CertificateAuthority.Data

	return nil
}
