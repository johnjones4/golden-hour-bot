tables:
	AWS_SECRET_ACCESS_KEY=dev AWS_ACCESS_KEY_ID=dev AWS_DEFAULT_REGION=us-east-1 AWS_PAGE="" aws dynamodb create-table --endpoint-url http://localhost:8000 --no-paginate --output=text --cli-input-yaml file://res/remindersTable.yml
	AWS_SECRET_ACCESS_KEY=dev AWS_ACCESS_KEY_ID=dev AWS_DEFAULT_REGION=us-east-1 AWS_PAGE="" aws dynamodb create-table --endpoint-url http://localhost:8000 --no-paginate --output=text --cli-input-yaml file://res/stateTable.yml
	AWS_SECRET_ACCESS_KEY=dev AWS_ACCESS_KEY_ID=dev AWS_DEFAULT_REGION=us-east-1 AWS_PAGE="" aws dynamodb create-table --endpoint-url http://localhost:8000 --no-paginate --output=text --cli-input-yaml file://res/remindersIndexTable.yml

drop-tables:
	AWS_SECRET_ACCESS_KEY=dev AWS_ACCESS_KEY_ID=dev AWS_DEFAULT_REGION=us-east-1 AWS_PAGE="" aws dynamodb delete-table --endpoint-url http://localhost:8000 --table-name golden-reminders
	AWS_SECRET_ACCESS_KEY=dev AWS_ACCESS_KEY_ID=dev AWS_DEFAULT_REGION=us-east-1 AWS_PAGE="" aws dynamodb delete-table --endpoint-url http://localhost:8000 --table-name golden-state
	AWS_SECRET_ACCESS_KEY=dev AWS_ACCESS_KEY_ID=dev AWS_DEFAULT_REGION=us-east-1 AWS_PAGE="" aws dynamodb delete-table --endpoint-url http://localhost:8000 --table-name golden-reminders-index

build:
	rm -rf bin || true
	cd runtimes/aws-alert-dequeuer && GOOS=linux go build -ldflags="-s -w" -o ../../bin/aws-alert-dequeuer
	cd runtimes/aws-alert-queuer && GOOS=linux go build -ldflags="-s -w" -o ../../bin/aws-alert-queuer
	cd runtimes/aws-webhook && GOOS=linux go build -ldflags="-s -w" -o ../../bin/aws-webhook
