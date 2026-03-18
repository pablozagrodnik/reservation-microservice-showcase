resource "digitalocean_database_cluster" "res_db" {
	name 		= "cin-res-db"
	engine 		= "pg"
	version 	= "15"
	size 		= "db-s-1vcpu-1gb"
	region 		= "fra1"
	node_count 	= 1
	private_network_uuid = digitalocean_vpc.res_vpc.id
}

resource "digitalocean_database_db" "res_db_name" {
	cluster_id = digitalocean_database_cluster.res_db.id
	name       = "cinebase"
}

