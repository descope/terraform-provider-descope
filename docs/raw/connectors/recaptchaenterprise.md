
RecaptchaEnterprise
===================



project_id
----------

- Type: `string` (required)

The Google Cloud project ID where the reCAPTCHA Enterprise is managed.



site_key
--------

- Type: `string` (required)

The site key is used to invoke reCAPTCHA Enterprise service on your site or
mobile application.



api_key
-------

- Type: `secret` (required)

API key associated with the current project.



override_assessment
-------------------

- Type: `bool` 

Override the default assessment model. Note: Overriding assessment is intended
for automated testing and should not be utilized in production environments.



assessment_score
----------------

- Type: `float` 
- Default: `0.5`

When configured, the Recaptcha action will return the score without assessing
the request. The score ranges between 0 and 1, where 1 is a human interaction
and 0 is a bot.
