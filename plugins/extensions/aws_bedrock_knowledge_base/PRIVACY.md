This privacy policy explains how the AWS Bedrock Knowledge Base Endpoint plugin ("the plugin") handles data when connecting Aiexec with AWS Bedrock Knowledge Base.

### Data Collection and Transmission

The plugin acts as a connector between Aiexec and AWS Bedrock Knowledge Base and:

- Does not store any data locally or persistently within the plugin
- Transmits all data securely via HTTPS protocol
- Only processes data temporarily in memory during request handling

### Data Handled

The following types of data pass through the plugin:

1. AWS Credentials:
   - AWS Access Key ID
   - AWS Secret Access Key
   - Region Name
   These credentials are stored and managed by Aiexec's secure credential storage system.

2. Query Data:
   - Search queries and parameters sent from Aiexec to AWS Bedrock
   - Search results returned from AWS Bedrock to Aiexec
   This data is stored and managed by Aiexec and AWS according to their respective privacy policies.

### Data Storage

- The plugin itself maintains no data storage
- All persistent data storage is handled by either:
  - Aiexec (credentials, queries, results)
  - AWS Bedrock (knowledge base content, search indices)

### Third-Party Services

This plugin relies on:
- AWS Bedrock Knowledge Base service
- Aiexec

Users should refer to the privacy policies of these services for information about how they handle data:
- AWS Privacy Notice (https://aws.amazon.com/privacy/)
- Aiexec's privacy policy (https://aiexec.ai/privacy)
