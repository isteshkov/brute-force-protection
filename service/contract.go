package service

import (
	"context"

	"gitlab.com/isteshkov/brute-force-protection/contract"
	myContext "gitlab.com/isteshkov/brute-force-protection/domain/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func newRPC(s *Service) *grpc.Server {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(s.MiddlewareAccess))
	contract.RegisterProtectorServer(grpcServer, s)
	reflection.Register(grpcServer)
	return grpcServer
}

func (s *Service) AuthAttempt(ctx context.Context,
	request *contract.RequestAuthAttempt) (response *contract.ResponseAuthAttempt, err error) {
	s.SetLogger(s.Logger.WithFields(myContext.LogFieldsFromGrpcContext(ctx)))
	response = &contract.ResponseAuthAttempt{}
	defer s.processRPCError(&err, &response.ErrorMsg)

	response.Allowed, err = s.authAttempt(request.Login, request.Password, request.IpAddress)
	if err != nil {
		return
	}

	return
}

func (s *Service) CleanBucketByLogin(ctx context.Context,
	request *contract.RequestCleanBucketByLogin) (response *contract.ResponseCleanBucketByLogin, err error) {
	s.SetLogger(s.Logger.WithFields(myContext.LogFieldsFromGrpcContext(ctx)))
	response = &contract.ResponseCleanBucketByLogin{}
	defer s.processRPCError(&err, &response.ErrorMsg)

	err = s.cleanBucketByLogin(request.Login)
	if err != nil {
		return
	}

	return
}

func (s *Service) CleanBucketByIP(ctx context.Context,
	request *contract.RequestCleanBucketByIP) (response *contract.ResponseCleanBucketByIP, err error) {
	s.SetLogger(s.Logger.WithFields(myContext.LogFieldsFromGrpcContext(ctx)))
	response = &contract.ResponseCleanBucketByIP{}
	defer s.processRPCError(&err, &response.ErrorMsg)

	err = s.cleanBucketByIP(request.IpAddress)
	if err != nil {
		return
	}

	return
}

func (s *Service) AddToBlackList(ctx context.Context,
	request *contract.RequestAddToList) (response *contract.ResponseAddToList, err error) {
	s.SetLogger(s.Logger.WithFields(myContext.LogFieldsFromGrpcContext(ctx)))
	response = &contract.ResponseAddToList{}
	defer s.processRPCError(&err, &response.ErrorMsg)

	err = s.addSubnetToBlacklist(request.SubnetAddress)
	if err != nil {
		return
	}

	return
}

func (s *Service) RemoveFromBlackList(ctx context.Context,
	request *contract.RequestRemoveFromList) (response *contract.ResponseRemoveFromList, err error) {
	s.SetLogger(s.Logger.WithFields(myContext.LogFieldsFromGrpcContext(ctx)))
	response = &contract.ResponseRemoveFromList{}
	defer s.processRPCError(&err, &response.ErrorMsg)

	err = s.removeSubnetFromList(request.SubnetAddress)
	if err != nil {
		return
	}

	return
}

func (s *Service) AddToWhiteList(ctx context.Context,
	request *contract.RequestAddToList) (response *contract.ResponseAddToList, err error) {
	s.SetLogger(s.Logger.WithFields(myContext.LogFieldsFromGrpcContext(ctx)))
	response = &contract.ResponseAddToList{}
	defer s.processRPCError(&err, &response.ErrorMsg)

	err = s.addSubnetToWhitelist(request.SubnetAddress)
	if err != nil {
		return
	}

	return
}

func (s *Service) RemoveFromWhiteList(ctx context.Context,
	request *contract.RequestRemoveFromList) (response *contract.ResponseRemoveFromList, err error) {
	s.SetLogger(s.Logger.WithFields(myContext.LogFieldsFromGrpcContext(ctx)))
	response = &contract.ResponseRemoveFromList{}
	defer s.processRPCError(&err, &response.ErrorMsg)

	err = s.removeSubnetFromList(request.SubnetAddress)
	if err != nil {
		return
	}

	return
}
