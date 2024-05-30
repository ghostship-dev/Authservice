CREATE MIGRATION m1bv7tti5k4fnb65ukkqt4lagnmcjve7jiy3vdnxaoeeqaj3chcyba
    ONTO m13zya64j63rrb5g5jykylxuku7ztn7cw32nhqiiatplx4fl552ffa
{
  ALTER TYPE default::Account {
      ALTER PROPERTY otp_state {
          CREATE CONSTRAINT std::one_of('disabled', 'enabled', 'verifying');
      };
  };
};
