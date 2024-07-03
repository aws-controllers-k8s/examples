{{- define "ack-pub-sub.topicPolicy"}}
{
  "Version": "2008-10-17",
  "Id": "__default_policy_ID",
  "Statement": [
    {
      "Sid": "__default_statement_ID",
      "Effect": "Allow",
      "Principal": {
        "AWS": "*"
      },
      "Action": [
        "SNS:GetTopicAttributes",
        "SNS:SetTopicAttributes",
        "SNS:AddPermission",
        "SNS:RemovePermission",
        "SNS:DeleteTopic",
        "SNS:Subscribe",
        "SNS:ListSubscriptionsByTopic",
        "SNS:Publish"
      ],
      "Resource": "arn:aws:sns:{{ .root.awsRegion }}:{{ .root.awsAccountNo }}:{{ .root.pubSubName }}-{{ .topic }}",
      "Condition": {
        "StringEquals": {
          "AWS:SourceOwner": "{{ .root.awsAccountNo }}"
        }
      }
    }
  ]
}
{{- end}}

{{- define "ack-pub-sub.queuePolicy"}}
{
  "Statement": [{
    "Sid": "__owner_statement",
    "Effect": "Allow",
    "Principal": {
      "AWS": "849707107029"
    },
    "Action": "sqs:SendMessage",
    "Resource": "arn:aws:sqs:{{ .root.awsRegion }}:{{ .root.awsAccountNo }}:{{ .root.pubSubName }}-{{ .topic }}-sqs-sub-{{ .topic }}"
  }]
}
{{- end}}