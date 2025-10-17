package verihubs

import (
	"backend-mobile-api/app/config"
	"backend-mobile-api/helpers"
	"backend-mobile-api/model/enum"
	verihubsDto "backend-mobile-api/model/outbond/verihubs-dto"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type outboundVeriHubsService struct {
	verihubsConfig *config.Verihubs
	rootConfig     *config.Root
	clog           *helpers.CustomLogger
}

func NewOutboundVeriHubsService(
	verihubsConfig *config.Verihubs,
	rootConfig *config.Root,
	clog *helpers.CustomLogger,
) OutboundVeriHubsService {
	return &outboundVeriHubsService{
		verihubsConfig: verihubsConfig,
		rootConfig:     rootConfig,
		clog:           clog,
	}

}

type OutboundVeriHubsService interface {
	VerifySMSOtpService(ctx context.Context, req *verihubsDto.VerifyOtpBaseRequest) (*verihubsDto.VerifyOtpBaseResponse, error)
	VerifyWhatsappsOtpService(ctx context.Context, req *verihubsDto.VerifyOtpBaseRequest) (*verihubsDto.VerifyOtpBaseResponse, error)
	verifyOtp(ctx context.Context, req *verihubsDto.VerifyOtpBaseRequest, url string, tag string) (*verihubsDto.VerifyOtpBaseResponse, error)
	SendWhatsappsService(ctx context.Context, req *verihubsDto.SendWhatsappOtpBaseRequest) (*verihubsDto.SendOtpWAResponse, error)
	SendSMSOtpService(ctx context.Context, req *verihubsDto.SendOtpBaseRequest) (*verihubsDto.SendOtpSMSResponse, error)
	SendKYCIdenityPassport(ctx context.Context, req *verihubsDto.VerihubIdentityRequest) (*verihubsDto.IdentityPassportResponse, error)
	SendKYCIdenityKTP(ctx context.Context, req *verihubsDto.VerihubIdentityRequest) (*verihubsDto.IdentityKTPResponse, error)
	SendVerifySelfie(ctx context.Context, req *verihubsDto.VerifyKycSelfie) (*verihubsDto.VerifySelfieResponse, error)
	SetHeaderRequest(header *http.Request) *http.Request
}

func (svc *outboundVeriHubsService) SetHeaderRequest(header *http.Request) *http.Request {
	header.Header.Add("Content-Type", "application/json")
	header.Header.Set("Accept", "application/json")
	header.Header.Set(string(enum.VERIHUBS_APP_ID), svc.verihubsConfig.AppID)
	header.Header.Set(string(enum.VERIHUBS_KEY), svc.verihubsConfig.VerihubsKey)
	return header
}

func (svc *outboundVeriHubsService) VerifySMSOtpService(ctx context.Context, req *verihubsDto.VerifyOtpBaseRequest) (*verihubsDto.VerifyOtpBaseResponse, error) {
	return svc.verifyOtp(ctx, req, fmt.Sprintf("%s%s", svc.verihubsConfig.Domain, "/v2/otp/verify"), "verifySMSOtpService")
}
func (svc *outboundVeriHubsService) VerifyWhatsappsOtpService(ctx context.Context, req *verihubsDto.VerifyOtpBaseRequest) (*verihubsDto.VerifyOtpBaseResponse, error) {
	return svc.verifyOtp(ctx, req, fmt.Sprintf("%s%s", svc.verihubsConfig.Domain, "/v1/whatsapp/otp/verify"), "VerifyWhatsappsOtpService")
}

func (svc *outboundVeriHubsService) verifyOtp(ctx context.Context, req *verihubsDto.VerifyOtpBaseRequest, url string, tag string) (*verihubsDto.VerifyOtpBaseResponse, error) {
	jsonReq, err := json.Marshal(*req)
	if err != nil {
		svc.clog.ErrorLogger(ctx, fmt.Sprintf("%s.%s", tag, "json.Marshal"), err)
		return nil, err
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
	if err != nil {
		svc.clog.ErrorLogger(ctx, fmt.Sprintf("%s.%s", tag, "http.NewRequest"), err)
		return nil, err
	}
	request = svc.SetHeaderRequest(request)
	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		svc.clog.ErrorLogger(ctx, fmt.Sprintf("%s.%s", tag, "httpClient.Do"), err)
		return nil, err
	}
	jsonBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		svc.clog.ErrorLogger(ctx, fmt.Sprintf("%s.%s", tag, "ioutil.ReadAll"), errors.New(http.StatusText(response.StatusCode)))
		if response.StatusCode != http.StatusCreated {
			return nil, errors.New(http.StatusText(response.StatusCode))
		}
		svc.clog.ErrorLogger(ctx, "sendOtp.ioutil.ReadAll", errors.New(string(jsonBody)))
		return nil, errors.New(string(jsonBody))
	}
	if response.StatusCode != http.StatusOK {
		svc.clog.ErrorLogger(ctx, fmt.Sprintf("%s.%s", tag, " statuscode != 201"), errors.New(string(jsonBody)))
		return nil, errors.New(string(jsonBody))
	}
	var responseBody verihubsDto.VerifyOtpBaseResponse
	if err = json.Unmarshal(jsonBody, &response); err != nil {
		svc.clog.ErrorLogger(ctx, fmt.Sprintf("%s.%s", tag, " json.Unmarshal"), errors.New(string(jsonBody)))
		return nil, errors.New(string(jsonBody))
	}
	return &responseBody, nil
}

func (svc *outboundVeriHubsService) SendWhatsappsService(ctx context.Context, req *verihubsDto.SendWhatsappOtpBaseRequest) (*verihubsDto.SendOtpWAResponse, error) {
	url := fmt.Sprintf("%s%s", svc.verihubsConfig.Domain, "/v1/whatsapp/otp/send")
	jsonReq, err := json.Marshal(*req)
	if err != nil {
		svc.clog.ErrorLogger(ctx, "SendWhatsappsService.json.Marshal", err)
		return nil, err
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
	if err != nil {
		svc.clog.ErrorLogger(ctx, "SendWhatsappsService.http.NewRequest", err)
		return nil, err
	}
	request = svc.SetHeaderRequest(request)
	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		svc.clog.ErrorLogger(ctx, "SendWhatsappsService.httpClient.Do", err)
		return nil, err
	}
	jsonBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		svc.clog.ErrorLogger(ctx, "SendWhatsappsService.ioutil.ReadAll", errors.New(http.StatusText(response.StatusCode)))
		if response.StatusCode != http.StatusCreated {
			return nil, errors.New(http.StatusText(response.StatusCode))
		}
		svc.clog.ErrorLogger(ctx, "sendOtp.ioutil.ReadAll", errors.New(string(jsonBody)))
		return nil, errors.New(string(jsonBody))
	}
	if response.StatusCode != http.StatusCreated {
		svc.clog.ErrorLogger(ctx, "SendWhatsappsService.statuscode != 201", errors.New(string(jsonBody)))
		return nil, errors.New(string(jsonBody))
	}
	var responseBody verihubsDto.SendOtpWAResponse
	if err = json.Unmarshal(jsonBody, &responseBody); err != nil {
		svc.clog.ErrorLogger(ctx, "SendWhatsappsService.json.Unmarshal", errors.New(string(jsonBody)))
		return nil, errors.New(string(jsonBody))
	}
	return &responseBody, nil
}

func (svc *outboundVeriHubsService) SendSMSOtpService(ctx context.Context, req *verihubsDto.SendOtpBaseRequest) (*verihubsDto.SendOtpSMSResponse, error) {
	url := fmt.Sprintf("%s%s", svc.verihubsConfig.Domain, "/v2/otp/send")
	jsonReq, err := json.Marshal(*req)

	if err != nil {
		svc.clog.ErrorLogger(ctx, "SendSMSOtpService.json.Marshal", err)
		return nil, err
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
	if err != nil {
		svc.clog.ErrorLogger(ctx, "SendSMSOtpService.http.NewRequest", err)
		return nil, err
	}
	request = svc.SetHeaderRequest(request)
	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		svc.clog.ErrorLogger(ctx, "SendSMSOtpService.httpClient.Do", err)
		return nil, err
	}
	jsonBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		svc.clog.ErrorLogger(ctx, "SendSMSOtpService.ioutil.ReadAll", errors.New(http.StatusText(response.StatusCode)))
		if response.StatusCode != http.StatusCreated {
			return nil, errors.New(http.StatusText(response.StatusCode))
		}
		svc.clog.ErrorLogger(ctx, "sendOtp.ioutil.ReadAll", errors.New(string(jsonBody)))
		return nil, errors.New(string(jsonBody))
	}
	if response.StatusCode != http.StatusCreated {
		svc.clog.ErrorLogger(ctx, "SendSMSOtpService.statuscode != 201", errors.New(string(jsonBody)))
		return nil, errors.New(string(jsonBody))
	}
	var responseBody verihubsDto.SendOtpSMSResponse
	if err = json.Unmarshal(jsonBody, &responseBody); err != nil {
		svc.clog.ErrorLogger(ctx, "SendSMSOtpService.json.Unmarshal", errors.New(string(jsonBody)))
		return nil, errors.New(string(jsonBody))
	}
	return &responseBody, nil

}

func (svc *outboundVeriHubsService) SendKYCIdenityKTP(ctx context.Context, req *verihubsDto.VerihubIdentityRequest) (*verihubsDto.IdentityKTPResponse, error) {
	var (
		result *verihubsDto.IdentityKTPResponse
		url    string
	)

	url = fmt.Sprintf("%s/v2/ktp/extract", svc.verihubsConfig.Domain)
	bodyRequest, err := json.Marshal(*req)
	if err != nil {
		svc.clog.ErrorLogger(ctx, "KYCIdenityKTP.statuscode != 201", errors.New(string(bodyRequest)))
		return nil, errors.New(string(bodyRequest))
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyRequest))
	if err != nil {
		svc.clog.ErrorLogger(ctx, "KYCIdenityKTP.httpNewRequest", err)
		return nil, err
	}

	request = svc.SetHeaderRequest(request)

	httpClient := &http.Client{}
	res, err := httpClient.Do(request)
	if err != nil {
		svc.clog.ErrorLogger(ctx, "KYCIdenityKTP.httpClient.Do", err)
		return nil, err
	}

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		svc.clog.ErrorLogger(ctx, "KYCIdenityKTP.ioutil.ReadAll", errors.New(http.StatusText(res.StatusCode)))
		if res.StatusCode != http.StatusCreated {
			return nil, errors.New(http.StatusText(res.StatusCode))
		}
		svc.clog.ErrorLogger(ctx, "KYCIdenityKTP.ioutil.ReadAll", errors.New(string(responseBody)))
		return nil, errors.New(string(responseBody))
	}

	if res.StatusCode != http.StatusOK {
		svc.clog.ErrorLogger(ctx, fmt.Sprintf("failed to e-kyc to verihubs for this request with response {%v}", responseBody), errors.New(string(responseBody)))
		return nil, errors.New(string(responseBody))
	}

	if err := json.Unmarshal(responseBody, &result); err != nil {
		svc.clog.ErrorLogger(ctx, fmt.Sprintf("failed unmarshal response body E-KYC KTP for res %v", res), errors.New(string(responseBody)))
		return nil, errors.New(string(responseBody))
	}

	return result, nil
}

func (svc *outboundVeriHubsService) SendKYCIdenityPassport(ctx context.Context, req *verihubsDto.VerihubIdentityRequest) (*verihubsDto.IdentityPassportResponse, error) {
	var (
		result *verihubsDto.IdentityPassportResponse
		url    string
	)
	url = fmt.Sprintf("%s/v2/ocr/passport", svc.verihubsConfig.Domain)

	// Decode base64 string
	decoded, err := base64.StdEncoding.DecodeString(req.Image)
	if err != nil {
		svc.clog.ErrorLogger(ctx, "KYCIdenityPassport.DecodeString", err)
		return nil, err
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", "verifyPasport/image")
	if err != nil {
		svc.clog.ErrorLogger(ctx, "KYCIdenityPassport.FailedToMultipartImage", err)
		return nil, err
	}

	_, err = io.Copy(part, bytes.NewReader(decoded))
	if err != nil {
		if err != nil {
			svc.clog.ErrorLogger(ctx, "KYCIdenityPassport.FailedToCopyImageFromPartToSrc", err)
			return nil, err
		}
	}

	writer.Close()

	reqVerihubs, err := http.NewRequest("POST", url, &body)
	if err != nil {
		svc.clog.ErrorLogger(ctx, "KYCPassport.http.NewRequest", err)
		return nil, err
	}

	reqVerihubs.Header.Set("Content-Type", writer.FormDataContentType())
	reqVerihubs.Header.Set(string(enum.VERIHUBS_APP_ID), svc.verihubsConfig.AppID)
	reqVerihubs.Header.Set(string(enum.VERIHUBS_KEY), svc.verihubsConfig.VerihubsKey)

	client := &http.Client{}
	resp, err := client.Do(reqVerihubs)
	if err != nil {
		svc.clog.ErrorLogger(ctx, "KYCPassport.httpClientDo", err)
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		svc.clog.ErrorLogger(ctx, "KYCPassport.ioutil.ReadAll", errors.New(http.StatusText(resp.StatusCode)))
		if resp.StatusCode != http.StatusCreated {
			return nil, errors.New(http.StatusText(resp.StatusCode))
		}
		svc.clog.ErrorLogger(ctx, "KYCPassport.ioutil.ReadAll", errors.New(string(responseBody)))
		return nil, errors.New(string(responseBody))
	}

	if resp.StatusCode != http.StatusOK {
		svc.clog.ErrorLogger(ctx, fmt.Sprintf("failed to e-kyc to verihubs for this request with response {%v}", responseBody), errors.New(string(responseBody)))
		return nil, errors.New(string(responseBody))
	}

	if err := json.Unmarshal(responseBody, &result); err != nil {
		svc.clog.ErrorLogger(ctx, fmt.Sprintf("failed unmarshal response body E-KYC KTP for res %v", resp), errors.New(string(responseBody)))
		return nil, errors.New(string(responseBody))
	}

	return result, nil
}

// SendVerifySelfie implements OutboundVeriHubsService.
func (svc *outboundVeriHubsService) SendVerifySelfie(ctx context.Context, req *verihubsDto.VerifyKycSelfie) (*verihubsDto.VerifySelfieResponse, error) {
	svc.clog.InfoLogger(ctx, "start send verify selfie")
	var (
		result *verihubsDto.VerifySelfieResponse
		url    string
	)

	url = fmt.Sprintf("%s/data-verification/certificate-electronic/verify", svc.verihubsConfig.Domain)
	bodyRequest, err := json.Marshal(*req)
	if err != nil {
		svc.clog.ErrorLogger(ctx, "SendVerifySelfie.statuscode != 201", errors.New(string(bodyRequest)))
		return nil, errors.New(string(bodyRequest))
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyRequest))
	if err != nil {
		svc.clog.ErrorLogger(ctx, "SendVerifySelfie.httpNewRequest", err)
		return nil, err
	}

	request = svc.SetHeaderRequest(request)

	httpClient := &http.Client{}
	res, err := httpClient.Do(request)
	if err != nil {
		svc.clog.ErrorLogger(ctx, "SendVerifySelfie.httpClient.Do", err)
		return nil, err
	}

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		svc.clog.ErrorLogger(ctx, "SendVerifySelfie.ioutil.ReadAll", errors.New(http.StatusText(res.StatusCode)))
		if res.StatusCode != http.StatusCreated {
			return nil, errors.New(http.StatusText(res.StatusCode))
		}
		svc.clog.ErrorLogger(ctx, "SendVerifySelfie.ioutil.ReadAll", errors.New(string(responseBody)))
		return nil, errors.New(string(responseBody))
	}

	if res.StatusCode != http.StatusOK {
		svc.clog.ErrorLogger(ctx, fmt.Sprintf("failed to e-kyc to verihubs for this request {%v} with response {%v}", req.Nik, responseBody), errors.New(string(responseBody)))
		var result verihubsDto.VerihubsErrorResponse
		if err := json.Unmarshal(responseBody, &result); err != nil {
			svc.clog.ErrorLogger(ctx, fmt.Sprintf("failed unmarshal error response body E-KYC Selfie for res %v", res), errors.New(string(responseBody)))
			return nil, errors.New(string(responseBody))
		}
		return nil, errors.New(result.ErrorFields[0].Message)
	}

	if err := json.Unmarshal(responseBody, &result); err != nil {
		svc.clog.ErrorLogger(ctx, fmt.Sprintf("failed unmarshal response body E-KYC Selfie for res %v", res), errors.New(string(responseBody)))
		return nil, errors.New(string(responseBody))
	}

	if result.Data.Status == "not_verified" {
		svc.clog.ErrorLogger(ctx, fmt.Sprintf("failed to verified selfie cause not verified field"), errors.New(string(responseBody)))
		return nil, fmt.Errorf("KYC failed '%s' is not verified or does not match", result.Data.RejectField[0])
	}
	svc.clog.InfoLogger(ctx, "Done send verify selfie")
	return result, nil
}
