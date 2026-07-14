\*\*\*Register flow:

- (auth/register/request-otp) input: phone -> server gen OTP code and cache in redis with expire time is 1m.
- (auth/register/verify-otp) input: phone, otp code -> server gen verification token and return.
- (auth/register/complete): input: verification token, full name, username, password -> server insert new user (decode verification token to obtain phone)
