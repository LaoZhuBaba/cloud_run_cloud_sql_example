# cloud_run_cloud_sql_example
An example demonstrating Google Cloud Run talking to Cloud SQL via VPC using Private Service Connect.

Also demonstrates how you can set a password safely using Terraform's ephemeral feature.  This allows you to create a password
with Terraform without exposing the password in the state file or code.  Also, you can rotate the secret by updating the variable: **local.password_version_number**.

The code can run in two different ways:  If you define the **DB_PASS** environment variable then the code will use **DB_USER** and **DB_PASS** for authentication.  But if DB_PASS is not provided then the golang code assumes that DB_USER is a Google IAM service account and it will query IAM for an Oauth2 token and use this as the password.  This approach is more modern and secure.  But there are few caveats: 

* The service account must be created as a user within Cloud SQL.  For some unknown reason, for Postgresql only when you create the DB login for the service account you have to use the service account's e-mail address __with the__ **.gserviceaccount.com** __suffix removed__.
* From experimentation I found that IAM accounts do not automatically have any permissions within Postgresql, so I found I had to do: **grant all privileges on schema public to "<ACCOUNT_NAME>"**;
* You must enable a feature flage on the SQL instance: **cloudsql.iam_authentication** in order to use IAM accounts
