package main

import "github.com/mdlayher/wifi"

type Service struct {
	client *wifi.Client
	iface  *wifi.Interface
}

func NewService(c *wifi.Client, i *wifi.Interface) (*Service, error) {
	return &Service{c, i}, nil
}

func (svc *Service) GetIface() *wifi.Interface {
	return svc.iface
}

func (svc *Service) GetStations() ([]*wifi.StationInfo, error) {
	return svc.client.StationInfo(svc.iface)
}
