script_path="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

aws iam delete-account-alias --account-alias pocok-local > /dev/null 2>&1
aws iam create-account-alias --account-alias pocok-local

account_id=`aws sts get-caller-identity --query "Account" --output text | sed 's/^ *//;s/ *$//'`

cp "$script_path/config.yml.example" "$script_path/config.yml"

sed -i '' "s/ACCOUNT_ID/$account_id/g" "$script_path/config.yml"