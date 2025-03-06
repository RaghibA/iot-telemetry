# Define environment variables
export POSTGRES_USER=admin
export POSTGRES_PASSWORD=password
export POSTGRES_DB=iot_telemetry
export POSTGRES_PORT=5432

export AUTH_HOST=auth-service
export AUTH_PORT=8080

export IOT_ADMIN_HOST=admin-service
export IOT_ADMIN_PORT=8081

export IOT_DATA_HOST=data-service
export IOT_DATA_PORT=8082

export CONSUMER_HOST=consumer-service
export CONSUMER_PORT=8083

export KAFKA_HOST=kafka
export KAFKA_PORT=9092

export JWT_SECRET=testJwtSecretKey

# Run Go tests recursively
.PHONY: test
test:
	@echo "Running tests with environment variables..."
	@env | grep -E 'POSTGRES_|AUTH_|IOT_|CONSUMER_|KAFKA_|JWT_SECRET'
	@go test -json ./... | jq -r ' \
	. | select(.Action=="fail" or .Action=="pass") | \
	if .Action == "fail" then \
			"[\u001b[31m\(.Action)\u001b[0m] \(.Package) \(.Test)" \
	else \
			"[\u001b[32m\(.Action)\u001b[0m] \(.Package) \(.Test)" \
	end \
	'

.PHONY: test-auth
test-auth:
	@echo "Running auth tests with environment variables..."
	@env | grep -E 'POSTGRES_|AUTH_|IOT_|CONSUMER_|KAFKA_|JWT_SECRET' # Debugging: print env vars
	@go test ./services/auth/... -v

.PHONY: test-admin
test-admin:
	@echo "Running admin tests with environment variables..."
	@env | grep -E 'POSTGRES_|AUTH_|IOT_|CONSUMER_|KAFKA_|JWT_SECRET' # Debugging: print env vars
	@go test ./services/admin/... -v