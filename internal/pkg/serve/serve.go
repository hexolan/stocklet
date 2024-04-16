// Copyright (C) 2024 Declan Teevan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package serve

// Port Definitions
const (
	grpcPort    string = "9090"
	gatewayPort string = "90"
)

// Get an address to a gRPC server using the standard port
func GetAddrToGrpc(host string) string {
	return host + ":" + grpcPort
}

// Get an address to a gRPC-gateway interface using the standard port
func GetAddrToGateway(host string) string {
	return host + ":" + gatewayPort
}
