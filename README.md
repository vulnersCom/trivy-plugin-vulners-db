# trivy-plugin-vulners-db

# To use vulners-db plugin, first you need to install it:
`trivy plugin install github.com/vulnersCom/trivy-plugin-vulners-db`

# Next step download vulners enriched database using vulners api-key
`trivy vulners-db --api-key <vulners api-key>`

# So, enjoy enriched database

For example try this out
`trivy image python:3.4-alpine -f json`