CREATE MIGRATION m1x5wdhi72up7gikfgvfidt22xxpps3aws7ofmlonbn2tfo2kevicq
    ONTO m1cjdgkihfmeimgnlqodkcwfs6af7hxmzb3gunlo5gzts4h7ma6l5q
{
  ALTER TYPE default::Account {
      CREATE INDEX ON (.email);
  };
  ALTER TYPE default::Password {
      CREATE INDEX ON (.email);
  };
};
