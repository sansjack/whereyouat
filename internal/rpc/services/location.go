package services

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"whereyouat/pkg/schemas/location"

	"github.com/oschwald/maxminddb-golang/v2"
)

type contextKey string

const RemoteAddrKey contextKey = "remoteAddr"

type LocationService struct {
	IpDb *maxminddb.Reader
	ctx  context.Context
}

func (ls *LocationService) SetContext(ctx context.Context) {
	ls.ctx = ctx
}

var record struct {
	Country struct {
		ISOCode string            `maxminddb:"iso_code"`
		Names   map[string]string `maxminddb:"names"`
	} `maxminddb:"country"`
}

func (ls *LocationService) Calculate(args *location.CalculateArgs, reply *location.CalculateReply) error {
	remoteAddr, ok := ls.ctx.Value(RemoteAddrKey).(string)
	if !ok {
		return fmt.Errorf("remote address not available")
	}

	ipStr, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		ipStr = remoteAddr
	}

	host, err := netip.ParseAddr(ipStr)
	if err != nil {
		return fmt.Errorf("failed to parse remote address: %w", err)
	}

	err = ls.IpDb.Lookup(host).Decode(&record)
	if err != nil {
		return fmt.Errorf("failed to lookup IP: %w", err)
	}

	*reply = location.CalculateReply{
		Location:   record.Country.Names["en"],
		//Confidence: 100, TODO: Implement confidence calculation based on IP data + latency
		IsoCode:    record.Country.ISOCode,
	}

	return nil
}
