package client

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

const (
	AWS_REGION    string = "eu-central-1"
	CLIENT_ID     string = ""
	CLIENT_SECRET string = ""
	RUNS_TABLE    string = "pipe_runs"
	ROLE          string = ""
)

type AwsClient struct {
	AwsRegion     string
	ClientId      string
	ClientSecret  string
	CognitoClient *cognitoidentityprovider.Client
	DynamoClient  *dynamodb.Client
	StsClient     *sts.Client
}

func NewAwsClient() CloudClient {

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile("steppes"))
	if err != nil {
		log.Fatalln(err)
	}
	cgclient := cognitoidentityprovider.NewFromConfig(cfg)
	stsclient := sts.NewFromConfig(cfg)
	cfg.Credentials = stscreds.NewAssumeRoleProvider(stsclient, ROLE)
	dbclient := dynamodb.NewFromConfig(cfg)
	return &AwsClient{
		AwsRegion:     AWS_REGION,
		ClientId:      CLIENT_ID,
		ClientSecret:  CLIENT_SECRET,
		CognitoClient: cgclient,
		DynamoClient:  dbclient,
		StsClient:     stsclient,
	}
}

func (awsclient *AwsClient) Login(username, password string) (*types.AuthenticationResultType, error) {
	ctx := context.Background()
	secret := awsclient.generateSecretHash(username)
	opts := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: &awsclient.ClientId,
		AuthParameters: map[string]string{
			"USERNAME":    username,
			"PASSWORD":    password,
			"SECRET_HASH": secret,
		},
	}
	authOut, err := awsclient.CognitoClient.InitiateAuth(ctx, opts)
	if err != nil {
		return nil, err
	}
	return authOut.AuthenticationResult, nil
}

func (awsclient *AwsClient) generateSecretHash(username string) string {
	hm := hmac.New(sha256.New, []byte(awsclient.ClientSecret))
	hm.Write([]byte(fmt.Sprintf("%s%s", username, awsclient.ClientId)))
	return base64.StdEncoding.EncodeToString(hm.Sum(nil))
}

func (awsclient *AwsClient) GetRuns() ([]PipelineRun, error) {
	var runs []PipelineRun
	ctx := context.Background()
	data, err := awsclient.DynamoClient.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(RUNS_TABLE),
		Limit:     aws.Int32(20),
	})
	if err != nil {
		return nil, err
	}
	log.Println(data.Items)
	err = attributevalue.UnmarshalListOfMaps(data.Items, &runs)
	if err != nil {
		return nil, err
	}
	return runs, nil
}
