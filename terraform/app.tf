resource "digitalocean_app" "res_service" {
	spec {
		name 	= "cin-res-api"
		region 	= "fra1"

		domain {
			name = "reservation.pawsdevhub.me"
			type = "PRIMARY"
		}

		alert {
			rule = "DEPLOYMENT_FAILED"
			destinations {
				emails = ["pawel.zagrodnik0@gmail.com"]
			}
		}

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
				value 	= "${digitalocean_database_cluster.res_db.private_uri}/cinebase"
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

		# Frontend
		static_site {
			name 			= "react-frontend"
			build_command 	= "npm run build"
			output_dir 		= "dist"
			source_dir 		= "frontend"

			git {
				repo_clone_url 	= "https://github.com/pablozagrodnik/reservation-microservice-showcase"
				branch 			= "main"
			}

			env {
				key = "VITE_API_URL"
				value = "https://reservation.pawsdevhub.me/api"
			}
		}

		ingress {
			rule {
				component {
					name = "go-backend"
				}
				match {
					path {
						prefix = "/api"
					}
				}
			}
			rule {
				component {
					name = "react-frontend"
				}
				match {
					path {
						prefix = "/"
					}
				}
			}
		}

		vpc {
			id = digitalocean_vpc.res_vpc.id
		}
	}

  lifecycle {
		ignore_changes = [
			spec[0].service[0].image[0].tag,
		]
  }
}
