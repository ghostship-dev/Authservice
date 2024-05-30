CREATE MIGRATION m1rwasruufpog5lh2lq62xgv2mwboqdakigbicmyxy4b3opnbhvnaa
    ONTO m1darjxigmqhll6eextmlsu24gv7cntumn6e63h6gmwfcfjj5lmkrq
{
  ALTER TYPE default::Authcode {
      CREATE REQUIRED PROPERTY redirect_uri: std::str {
          SET REQUIRED USING (<std::str>{});
      };
  };
};
