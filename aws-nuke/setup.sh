script_path="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

aws iam create-account-alias --account-alias local-test

account_id=`aws sts get-caller-identity --query "Account" --output text | sed 's/^ *//;s/ *$//'`

cp "$script_path/config.yml.example" "$script_path/config.yml"

sed -i '' "s/ACCOUNT_ID/$account_id/g" "$script_path/config.yml"