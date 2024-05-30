CREATE MIGRATION m1l2ryn5vglpivib2wpdix6ofogwjck65unb6t5varn4fiwdw44boq
    ONTO m1bv7tti5k4fnb65ukkqt4lagnmcjve7jiy3vdnxaoeeqaj3chcyba
{
  ALTER TYPE default::Token {
      ALTER PROPERTY revoked {
          SET REQUIRED USING (<std::bool>{});
      };
  };
};
