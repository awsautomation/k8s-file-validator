# Configure the AWS provider with appropriate credentials and region
provider "aws" {
  profile = "production"
  region  = "us-west-2"
}

# Create the IAM role for cross-account access
resource "aws_iam_role" "cross_account_role" {
  name = "instaclustr"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::31111111114357:role/root"
      },
      "Action": "sts:AssumeRole",
      "Condition": {
        "StringEquals": {
          "sts:ExternalId": "YOUR_INSTACLUSTR_CONSOLE_ACCOUNT_ID"
        }
      }
    }
  ]
}
EOF
}

# Attach the policy to the cross-account role
resource "aws_iam_role_policy_attachment" "cross_account_policy_attachment" {
  role       = aws_iam_role.cross_account_role.name
  policy_arn = "ARN_OF_INSTACLUSTR_POLICY"
}

# Modify the trust relationship for the cross-account role
resource "aws_iam_role" "cross_account_role_edit" {
  name               = "instaclustr"
  assume_role_policy = aws_iam_role.cross_account_role.assume_role_policy
}

data "aws_iam_policy_document" "cross_account_role_trust_policy" {
  statement {
    effect = "Allow"
    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::340600000057:role/InstaclustrProvisioning"]
    }
    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role_policy" "cross_account_role_trust_policy" {
  name   = "cross-account-role-trust-policy"
  role   = aws_iam_role.cross_account_role_edit.name
  policy = data.aws_iam_policy_document.cross_account_role_trust_policy.json
}
