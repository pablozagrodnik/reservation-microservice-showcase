resource "digitalocean_container_registry" "main" {
	name 					= "pawsdevhub-registry"
	subscription_tier_slug 	= "starter"
	region 					= "fra1"
}
