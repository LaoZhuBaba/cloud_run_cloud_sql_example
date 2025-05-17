# cloud_run_cloud_sql_example
An example demonstrating Google Cloud Run talking to Cloud SQL via VPC using Private Service Connect.

Also demostrates how you can set a password safely using Terraform's ephemeral feature.  This allows you to create a password
with Terraform without exposing the password in the state file or code.
