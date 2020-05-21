# Cake Team Survey

A serverless survey runner for Cake Experiment inception surveys

## Deploying

Built with the [serverless framework](https://serverless.com)

```shell
npm install

export <AWS_CREDENTIALS>
make deploy
```

## Resources created

- Survey lambda - main survey runner and end user interaction
- Admin lambda - back office functions
- DynamoDB table - holds all persisted data in single-table design

## Plugins

- [serverless-s3-sync](https://www.npmjs.com/package/serverless-s3-sync)
