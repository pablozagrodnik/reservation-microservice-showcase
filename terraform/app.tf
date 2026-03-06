resource "digitalocean_app" "res_service" {
	spec {
		name 	= "cin-res-api"
		region 	= "fra1"

		service {
			name 				= "go-backend"
			instance_count 		= 1
			http_port			= 8080
			instance_size_slug 	= "apps-s-1vcpu-0.5gb"

			image {
				registry_type	= "DOCR"
				repository		= "cinema-backend"
				tag				= var.image_tag
			}

			env {
				key 	= "DATABASE_URL"
				value 	= digitalocean_database_cluster.res_db.private_uri
				type 	= "SECRET"
			}

			env {
				key 	= "PORT"
				value 	= "8080"
			}

			health_check {
				http_path				= "/"
				initial_delay_seconds	= 10
				period_seconds			= 30
				timeout_seconds			= 5
				success_threshold		= 1
				failure_threshold		= 3
			}
		}
	}
}
