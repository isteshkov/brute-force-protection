package service

import (
	"context"

	"gitlab.com/isteshkov/brute-force-protection/contract"
	myContext "gitlab.com/isteshkov/brute-force-protection/domain/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func newRpc(s *Service) *grpc.Server {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(s.MiddlewareAccess))
	contract.RegisterProtectorServer(grpcServer, s)
	reflection.Register(grpcServer)
	return grpcServer
}

func (s *Service) AuthAttempt(ctx context.Context, request *contract.RequestAuthAttempt) (response *contract.ResponseAuthAttempt, err error) {
	s.SetLogger(s.Logger.WithFields(myContext.LogFieldsFromGrpcContext(ctx)))
	response = &contract.ResponseAuthAttempt{}
	defer s.processRpcError(&err, &response.ErrorMsg)

	err = s.authAttempt(request.Login, request.Password, request.IpAddress)
	if err != nil {
		return
	}

	return
}

func (s *Service) CleanBucketByLogin(ctx context.Context, request *contract.RequestCleanBucketByLogin) (response *contract.ResponseCleanBucketByLogin, err error) {
	s.SetLogger(s.Logger.WithFields(myContext.LogFieldsFromGrpcContext(ctx)))
	response = &contract.ResponseCleanBucketByLogin{}
	defer s.processRpcError(&err, &response.ErrorMsg)

	err = s.cleanBucketByLogin(request.Login)
	if err != nil {
		return
	}

	return
}

func (s *Service) CleanBucketByIp(ctx context.Context, request *contract.RequestCleanBucketByIp) (response *contract.ResponseCleanBucketByIp, err error) {
	s.SetLogger(s.Logger.WithFields(myContext.LogFieldsFromGrpcContext(ctx)))
	response = &contract.ResponseCleanBucketByIp{}
	defer s.processRpcError(&err, &response.ErrorMsg)

	err = s.cleanBucketByIp(request.IpAddress)
	if err != nil {
		return
	}

	return
}

func (s *Service) AddToBlackList(ctx context.Context, request *contract.RequestAddToList) (response *contract.ResponseAddToList, err error) {
	s.SetLogger(s.Logger.WithFields(myContext.LogFieldsFromGrpcContext(ctx)))
	response = &contract.ResponseAddToList{}
	defer s.processRpcError(&err, &response.ErrorMsg)

	err = s.addSubnetToBlacklist(request.SubnetAddress)
	if err != nil {
		return
	}

	return
}

func (s *Service) RemoveFromBlackList(ctx context.Context, request *contract.RequestRemoveFromList) (response *contract.ResponseRemoveFromList, err error) {
	s.SetLogger(s.Logger.WithFields(myContext.LogFieldsFromGrpcContext(ctx)))
	response = &contract.ResponseRemoveFromList{}
	defer s.processRpcError(&err, &response.ErrorMsg)

	err = s.removeSubnetFromBlacklist(request.SubnetAddress)
	if err != nil {
		return
	}

	return
}

func (s *Service) AddToWhiteList(ctx context.Context, request *contract.RequestAddToList) (response *contract.ResponseAddToList, err error) {
	s.SetLogger(s.Logger.WithFields(myContext.LogFieldsFromGrpcContext(ctx)))
	response = &contract.ResponseAddToList{}
	defer s.processRpcError(&err, &response.ErrorMsg)

	err = s.addSubnetToWhitelist(request.SubnetAddress)
	if err != nil {
		return
	}

	return
}

func (s *Service) RemoveFromWhiteList(ctx context.Context, request *contract.RequestRemoveFromList) (response *contract.ResponseRemoveFromList, err error) {
	s.SetLogger(s.Logger.WithFields(myContext.LogFieldsFromGrpcContext(ctx)))
	response = &contract.ResponseRemoveFromList{}
	defer s.processRpcError(&err, &response.ErrorMsg)

	err = s.removeSubnetFromWhitelist(request.SubnetAddress)
	if err != nil {
		return
	}

	return
}