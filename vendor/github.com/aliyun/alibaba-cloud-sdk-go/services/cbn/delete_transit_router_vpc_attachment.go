package cbn

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// DeleteTransitRouterVpcAttachment invokes the cbn.DeleteTransitRouterVpcAttachment API synchronously
func (client *Client) DeleteTransitRouterVpcAttachment(request *DeleteTransitRouterVpcAttachmentRequest) (response *DeleteTransitRouterVpcAttachmentResponse, err error) {
	response = CreateDeleteTransitRouterVpcAttachmentResponse()
	err = client.DoAction(request, response)
	return
}

// DeleteTransitRouterVpcAttachmentWithChan invokes the cbn.DeleteTransitRouterVpcAttachment API asynchronously
func (client *Client) DeleteTransitRouterVpcAttachmentWithChan(request *DeleteTransitRouterVpcAttachmentRequest) (<-chan *DeleteTransitRouterVpcAttachmentResponse, <-chan error) {
	responseChan := make(chan *DeleteTransitRouterVpcAttachmentResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeleteTransitRouterVpcAttachment(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// DeleteTransitRouterVpcAttachmentWithCallback invokes the cbn.DeleteTransitRouterVpcAttachment API asynchronously
func (client *Client) DeleteTransitRouterVpcAttachmentWithCallback(request *DeleteTransitRouterVpcAttachmentRequest, callback func(response *DeleteTransitRouterVpcAttachmentResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeleteTransitRouterVpcAttachmentResponse
		var err error
		defer close(result)
		response, err = client.DeleteTransitRouterVpcAttachment(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// DeleteTransitRouterVpcAttachmentRequest is the request struct for api DeleteTransitRouterVpcAttachment
type DeleteTransitRouterVpcAttachmentRequest struct {
	*requests.RpcRequest
	ResourceOwnerId           requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ClientToken               string           `position:"Query" name:"ClientToken"`
	DryRun                    requests.Boolean `position:"Query" name:"DryRun"`
	ResourceOwnerAccount      string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount              string           `position:"Query" name:"OwnerAccount"`
	OwnerId                   requests.Integer `position:"Query" name:"OwnerId"`
	ResourceType              string           `position:"Query" name:"ResourceType"`
	TransitRouterAttachmentId string           `position:"Query" name:"TransitRouterAttachmentId"`
}

// DeleteTransitRouterVpcAttachmentResponse is the response struct for api DeleteTransitRouterVpcAttachment
type DeleteTransitRouterVpcAttachmentResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateDeleteTransitRouterVpcAttachmentRequest creates a request to invoke DeleteTransitRouterVpcAttachment API
func CreateDeleteTransitRouterVpcAttachmentRequest() (request *DeleteTransitRouterVpcAttachmentRequest) {
	request = &DeleteTransitRouterVpcAttachmentRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cbn", "2017-09-12", "DeleteTransitRouterVpcAttachment", "", "")
	request.Method = requests.POST
	return
}

// CreateDeleteTransitRouterVpcAttachmentResponse creates a response to parse from DeleteTransitRouterVpcAttachment response
func CreateDeleteTransitRouterVpcAttachmentResponse() (response *DeleteTransitRouterVpcAttachmentResponse) {
	response = &DeleteTransitRouterVpcAttachmentResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}