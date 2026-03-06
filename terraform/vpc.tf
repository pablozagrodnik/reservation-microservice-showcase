resource "digitalocean_vpc" "res_vpc" {
	name 		= "cin-res-net"
	region 		= "fra1"
	ip_range 	= "10.10.10.0/24"
}
