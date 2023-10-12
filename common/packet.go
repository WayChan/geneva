package common

import (
	"encoding/binary"

	"github.com/google/gopacket/layers"
)

// UpdateTCPChecksum updates the TCP checksum field and the raw bytes for a gopacket TCP layer.
func UpdateTCPChecksum(tcp *layers.TCP) error {
	// the ComputeChecksum method requires the checksum bytes in the raw packet to be zeroed out.
	tcp.Contents[16] = 0
	tcp.Contents[17] = 0

	chksum, err := tcp.ComputeChecksum()
	if err != nil {
		return err
	}

	tcp.Checksum = chksum
	binary.BigEndian.PutUint16(tcp.Contents[16:18], chksum)

	return nil
}

// UpdateIPv4Checksum updates the IPv4 checksum field and the raw bytes for a gopacket IPv4 layer.
func UpdateIPv4Checksum(ip *layers.IPv4) error {
	buf := make([]byte, ip.Length)
	copy(buf, ip.Contents)
	copy(buf[len(ip.Contents):], ip.Payload)

	chksum := CalculateIPv4Checksum(buf)
	ip.Checksum = chksum
	binary.BigEndian.PutUint16(ip.Contents[10:12], chksum)

	return nil
}

// copied directly from gopacket/layers/ip4.go because they didn't export one. for whatever some reason..
func CalculateIPv4Checksum(bytes []byte) uint16 {
	// Clear checksum bytes
	bytes[10] = 0
	bytes[11] = 0

	// Compute checksum
	var csum uint32
	for i := 0; i < len(bytes); i += 2 {
		csum += uint32(bytes[i]) << 8
		csum += uint32(bytes[i+1])
	}

	for csum > 0xFFFF {
		// Add carry to the sum
		csum = (csum >> 16) + uint32(uint16(csum))
	}
	// Flip all the bits
	return ^uint16(csum)
}
