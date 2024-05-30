CREATE MIGRATION m1ztcgm6vgygcxx5x2jrhbvojrduqlkngqyntlakup7z3riuwl4rpa
    ONTO m1x5wdhi72up7gikfgvfidt22xxpps3aws7ofmlonbn2tfo2kevicq
{
  ALTER TYPE default::Token {
      CREATE INDEX ON (.token_value);
  };
};
