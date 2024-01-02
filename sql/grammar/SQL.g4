// Define a grammar called LinSQL for LinDB query language
// antlr4 SQL.g4 -Dlanguage=Go -package grammar
grammar SQL;

statement               : showStmt
                        | createBrokerStmt
                        | recoverStorageStmt
                        | useStmt
                        | queryStmt
                        | createDatabaseStmt
                        | dropDatabaseStmt
                        | setLimitStmt
                        | ident // just for suggest filtering.
                        EOF ;

useStmt                 : T_USE ident ;
setLimitStmt            : T_SET T_LIMIT toml;

showStmt                : showMasterStmt
                        | showMetadataTypesStmt
                        | showMetadatasStmt
                        | showBrokersStmt
                        | showLimitStmt
                        | showAliveStmt
                        | showStateStmt
                        | showReplicationsStmt
                        | showDatabasesStmt
                        | showRequestsStmt
                        // show metric metadata
                        | showNamespacesStmt
                        | showMetricsStmt
                        | showFieldsStmt
                        | showTagKeysStmt
                        | showTagValuesStmt
                        ;
//meta data query statement
showMasterStmt       	: T_SHOW T_MASTER ;
showBrokersStmt      	: T_SHOW T_BROKERS ;
showRequestsStmt     	: T_SHOW T_REQUESTS (T_WHERE T_ID T_EQUAL requestID)?; 
recoverStorageStmt   	: T_RECOVER T_STORAGE;
showLimitStmt        	: T_SHOW T_LIMIT ; // database limit
showMetadataTypesStmt	: T_SHOW T_METADATA T_TYPES;
showMetadatasStmt    	: T_SHOW T_METADATAS whereClause?;
showAliveStmt        	: T_SHOW (T_ROOT | T_BROKER | T_STORAGE) T_ALIVE;
showReplicationsStmt 	: T_SHOW T_REPLICATIONS whereClause?;
showStateStmt        	: T_SHOW T_STATE whereClause?;

// database
createDatabaseStmt   	: T_CREATE T_DATASBAE name (T_WITH properties)? (T_ROLLUP properties)?;
showDatabasesStmt    	: T_SHOW T_DATASBAES whereClause?;
dropDatabaseStmt     	: T_DROP T_DATASBAE name;
flushDatabaseStmt    	: T_FLUSH T_DATASBAE name;
compactDatabaseStmt  	: T_COMPACT T_DATASBAE name;

createBrokerStmt     	: T_CREATE T_BROKER name (T_WITH properties)?;
showNamespacesStmt   	: T_SHOW T_NAMESPACES (T_WHERE T_NAMESPACE T_EQUAL prefix)? limitClause?;
showMetricsStmt      	: T_SHOW T_METRICS (T_ON namespace)? (T_WHERE T_METRIC T_EQUAL prefix)? limitClause?;
showFieldsStmt       	: T_SHOW T_FIELDS fromClause;
showTagKeysStmt      	: T_SHOW T_TAG T_KEYS fromClause;
showTagValuesStmt    	: T_SHOW T_TAG T_VALUES fromClause T_WITH T_KEY T_EQUAL withTagKey whereClause? limitClause?;

//data query plan
queryStmt            	: T_EXPLAIN? sourceAndSelect whereClause? groupByClause? orderByClause? limitClause?;
sourceAndSelect      	: selectExpr fromClause | fromClause selectExpr ;
selectExpr           	: T_SELECT fields;
//select fields
fields               	: field ( T_COMMA field )* ;
field                	: fieldExpr alias? ;
alias                	: T_AS ident ;

fromClause           	: T_FROM metricName (T_ON namespace)? ;                 // from clause
whereClause          	: T_WHERE conditionExpr;                                // where clause
conditionExpr        	: expression ;
expression              : timeExpr                                              // time range
                        | name T_EQUAL value                                    // =
                        | name (T_NOTEQUAL | T_NOTEQUAL2) value                 // !=
                        | name T_NOT? T_LIKE value                              // like
                        | name T_REGEXP value                                   // =~: matches against 
                        | name T_NEQREGEXP value                                // !~: doesn’t match against
                        | name T_NOT? T_IN T_OPEN_P valueList T_CLOSE_P         // in
                        | expression T_AND expression                           // and
                        | expression T_OR expression                            // or
                        | T_OPEN_P expression T_CLOSE_P
                        ;


valueList            	: value (T_COMMA value)*;
timeExpr             	: T_TIME binaryOperator (nowExpr | ident) ;
nowExpr                 : nowFunc  durationLit? ;
nowFunc                 : T_NOW T_OPEN_P exprFuncParams? T_CLOSE_P ;

//group by
groupByClause           : T_GROUP T_BY groupByKeys (T_FILL T_OPEN_P fillOption T_CLOSE_P)? havingClause? ;
groupByKeys             : groupByKey (T_COMMA groupByKey)* ;
groupByKey              : ident | T_TIME T_OPEN_P durationLit T_CLOSE_P | T_TIME T_OPEN_P T_CLOSE_P;
fillOption              : T_NULL | T_PREVIOUS | L_INT | L_DEC ;

orderByClause           : T_ORDER T_BY sortFields ;
sortField               : fieldExpr ( T_ASC | T_DESC )* ;
sortFields              : sortField ( T_COMMA sortField )* ;

havingClause            : T_HAVING boolExpr ;
boolExpr                : T_OPEN_P boolExpr T_CLOSE_P
                        | boolExpr boolExprLogicalOp boolExpr
                        | boolExprAtom
                        ;
boolExprLogicalOp     	: T_AND  | T_OR ;
boolExprAtom            : binaryExpr ;
binaryExpr              : fieldExpr binaryOperator fieldExpr;
binaryOperator          : T_EQUAL
                        | T_NOTEQUAL
                        | T_NOTEQUAL2
                        | T_LESS
                        | T_LESSEQUAL
                        | T_GREATER
                        | T_GREATEREQUAL
                        | (T_LIKE | T_REGEXP)
                        ;

fieldExpr               : fieldExpr T_MUL fieldExpr
                        | fieldExpr T_DIV fieldExpr
                        | fieldExpr T_ADD fieldExpr
                        | fieldExpr T_SUB fieldExpr
                        | T_OPEN_P fieldExpr T_CLOSE_P
                        | exprFunc
                        | exprAtom
                        | durationLit
                        | star
                        ;
star                    : T_MUL;

durationLit             : intNumber intervalItem ;
intervalItem            : T_SECOND
                        | T_MINUTE
                        | T_HOUR
                        | T_DAY
                        | T_WEEK
                        | T_MONTH
                        | T_YEAR
                        ;
exprFunc                : funcName T_OPEN_P exprFuncParams? T_CLOSE_P ;
funcName                : T_SUM | T_MIN | T_MAX | T_AVG | T_COUNT | T_LAST | T_FIRST | T_STDDEV | T_QUANTILE | T_RATE;
exprFuncParams          : funcParam (T_COMMA funcParam)* ;
funcParam               : fieldExpr ;
exprAtom                : ident 
                        | decNumber
                        | intNumber
                        ;

properties              : T_OPEN_P propertyAssignments T_CLOSE_P ;
propertyAssignments     : property (T_COMMA property)* ;
property                : ident T_COLON value ;

value                   : STRING
                        | intNumber
                        | decNumber
                        | 'true'
                        | 'false'
                        | 'null'
                        | durationLit
                        ;

// Integer (positive or negative)
intNumber               : ('-' | '+')? L_INT ;
// Decimal number (positive or negative)
decNumber               : ('-' | '+')? L_DEC ;
limitClause             : T_LIMIT L_INT ;

metricName              : ident ;
tagKey                  : ident ;
tagValue                : ident ;
prefix                  : ident ;
withTagKey              : ident ;
namespace               : ident ;
name                    : ident ;
requestID               : ident ;
toml                    : ident ;

ident                   : (L_ID | nonReservedWords) ('.' (L_ID | nonReservedWords))* ;

nonReservedWords        : T_CREATE
                        | T_UPDATE
                        | T_SET
                        | T_DROP
                        | T_FLUSH
                        | T_COMPACT
                        | T_INTERVAL
                        | T_INTERVAL_NAME
                        | T_SHARD
                        | T_REPLICATIONS
                        | T_MEMORY
                        | T_TTL
                        | T_META_TTL
                        | T_PAST_TTL
                        | T_FUTURE_TTL
                        | T_KILL
                        | T_ON
                        | T_SHOW
                        | T_DATASBAE
                        | T_DATASBAES
                        | T_NAMESPACE
                        | T_NAMESPACES
                        | T_NODE
                        | T_METRICS
                        | T_METRIC
                        | T_FIELD
                        | T_FIELDS
                        | T_TAG
                        | T_INFO
                        | T_KEYS
                        | T_KEY
                        | T_WITH
                        | T_VALUES
                        | T_VALUE
                        | T_FROM
                        | T_WHERE
                        | T_LIMIT
                        | T_QUERIES
                        | T_QUERY
                        | T_EXPLAIN
                        | T_SELECT
                        | T_AS
                        | T_AND
                        | T_OR
                        | T_FILL
                        | T_NULL
                        | T_PREVIOUS
                        | T_ORDER
                        | T_ASC
                        | T_DESC
                        | T_LIKE
                        | T_NOT
                        | T_BETWEEN
                        | T_IS
                        | T_GROUP
                        | T_HAVING
                        | T_BY
                        | T_FOR
                        | T_STATS
                        | T_TIME
                        | T_NOW
                        | T_IN
                        | T_LOG
                        | T_PROFILE
                        | T_SUM
                        | T_MIN
                        | T_MAX
                        | T_COUNT
                        | T_LAST
                        | T_FIRST
                        | T_AVG
                        | T_STDDEV
                        | T_QUANTILE
                        | T_RATE
                        | T_SECOND
                        | T_MINUTE
                        | T_HOUR
                        | T_DAY
                        | T_WEEK
                        | T_MONTH
                        | T_YEAR
                        | T_USE
                        | T_MASTER
                        | T_METADATA
                        | T_METADATAS
                        | T_TYPE
                        | T_TYPES
                        | T_STORAGE
                        | T_ALIVE
                        | T_BROKER
                        | T_ROOT
                        | T_BROKERS
                        | T_SCHEMAS
                        | T_STATE_REPO
                        | T_STATE_MACHINE
                        | T_REQUESTS
                        | T_REQUEST
                        | T_ID
                        | T_ROLLUP
                        ;

STRING                  : '"' (ESC | SAFECODEPOINT)* '"' ;

fragment ESC
   : '\\' (["\\/bfnrt] | UNICODE)
   ;
fragment UNICODE
   : 'u' HEX HEX HEX HEX
   ;
fragment HEX
   : [0-9a-fA-F]
   ;
fragment SAFECODEPOINT
   : ~ ["\\\u0000-\u001F]
   ;

// no leading zeros

fragment EXP
   : [Ee] [+\-]? L_INT
   ;

// \- since - means "range" inside [...]
WS
   : [ \t\n\r] + -> skip
   ;

// Lexer rules
T_CREATE             : C R E A T E                      ;
T_UPDATE             : U P D A T E                      ;
T_SET                : S E T                            ;
T_DROP               : D R O P                          ;
T_FLUSH              : F L U S H                        ;
T_COMPACT            : C O M P A C T                    ;
T_INTERVAL           : I N T E R V A L                  ;
T_INTERVAL_NAME      : N A M E                          ;
T_SHARD              : S H A R D                        ;
T_REPLICATIONS       : R E P L I C A T I O N S          ;
T_MEMORY             : M E M O R Y                      ;
T_TTL                : T T L                            ;
T_META_TTL           : M E T A T T L                    ;
T_PAST_TTL           : P A S T T T L                    ;
T_FUTURE_TTL         : F U T U R E T T L                ;
T_KILL               : K I L L                          ;
T_ON                 : O N                              ;
T_SHOW               : S H O W                          ;
T_RECOVER            : R E C O V E R                    ;
T_USE                : U S E                            ;
T_STATE              : S T A T E                        ;
T_STATE_REPO         : S T A T E T_UNDERLINE R E P O    ;
T_STATE_MACHINE      : S T A T E T_UNDERLINE M A C H I N E;
T_MASTER             : M A S T E R                      ;
T_METADATA           : M E T A D A T A                  ;
T_METADATAS          : M E T A D A T A S                ;
T_TYPES              : T Y P E S                        ;
T_TYPE               : T Y P E                          ;
T_STORAGE            : S T O R A G E                    ;
T_BROKER             : B R O K E R                      ;
T_ROOT               : R O O T                          ;
T_BROKERS            : B R O K E R S                    ;
T_ALIVE              : A L I V E                        ;
T_SCHEMAS            : S C H E M A S                    ;
T_DATASBAE           : D A T A B A S E                  ;
T_DATASBAES          : D A T A B A S E S                ;
T_NAMESPACE          : N A M E S P A C E                ;
T_NAMESPACES         : N A M E S P A C E S              ;
T_NODE               : N O D E                          ;
T_METRICS            : M E T R I C S                    ;
T_METRIC             : M E T R I C                      ;
T_FIELD              : F I E L D                        ;
T_FIELDS             : F I E L D S                      ;
T_TAG                : T A G                            ;
T_INFO               : I N F O                          ;
T_KEYS               : K E Y S                          ;
T_KEY                : K E Y                            ;
T_WITH               : W I T H                          ;
T_VALUES             : V A L U E S                      ;
T_VALUE              : V A L U E                        ;
T_FROM               : F R O M                          ;
T_WHERE              : W H E R E                        ;
T_LIMIT              : L I M I T                        ;
T_QUERIES            : Q U E R I E S                    ;
T_QUERY              : Q U E R Y                        ;
T_EXPLAIN            : E X P L A I N                    ;
T_SELECT             : S E L E C T                      ;
T_AS                 : A S                              ;
T_AND                : A N D                            ;
T_OR                 : O R                              ;
T_FILL               : F I L L                          ;
T_NULL               : N U L L                          ;
T_PREVIOUS           : P R E V I O U S                  ;
T_ORDER              : O R D E R                        ;
T_ASC                : A S C                            ;
T_DESC               : D E S C                          ;
T_LIKE               : L I K E                          ;
T_NOT                : N O T                            ;
T_BETWEEN            : B E T W E E N                    ;
T_IS                 : I S                              ;
T_GROUP              : G R O U P                        ;
T_HAVING             : H A V I N G                      ;
T_BY                 : B Y                              ;
T_FOR                : F O R                            ;
T_STATS              : S T A T S                        ;
T_TIME               : T I M E                          ;
T_NOW                : N O W                            ;
T_IN                 : I N                              ;
T_ROLLUP             : R O L L U P                      ;

T_LOG                : L O G                            ;
T_PROFILE            : P R O F I L E                    ;
T_REQUESTS           : R E Q U E S T S                  ;
T_REQUEST            : R E Q U E S T                    ;
T_ID                 : I D                              ;

T_SUM                : S U M                            ;
T_MIN                : M I N                            ;
T_MAX                : M A X                            ;
T_COUNT              : C O U N T                        ;
T_LAST               : L A S T                          ;
T_FIRST              : F I R S T                        ;
T_AVG                : A V G                            ;
T_STDDEV             : S T D D E V                      ;
T_QUANTILE           : Q U A N T I L E                  ;
T_RATE               : R A T E                          ;

// create table option key
T_NUM_OF_SHARD   : N U M O F S H A R D;
T_REPLICA_FACTOR : R E P L I C A F A C T O R;
T_AUTO_CREATE_NS : A U T O C R E A T E N S;
T_BEHEAD         : B E H E A D;
T_AHEAD          : A H E A D;
T_RETENTION      : R E T E N T I O N;

//time unit
T_SECOND             : S                                ;
T_MINUTE             : 'm'                              ;
T_HOUR               : H                                ;
T_DAY                : D                                ;
T_WEEK               : W                                ;
T_MONTH              : 'M'                              ;
T_YEAR               : Y                                ;

//
T_DOT                :  '.'   ;
T_COLON              :  ':'   ;
T_EQUAL              :  '='   ;
T_NOTEQUAL           :  '<>'  ;
T_NOTEQUAL2          :  '!='  ;
T_GREATER            :  '>'   ;
T_GREATEREQUAL       :  '>='  ;
T_LESS               :  '<'   ;
T_LESSEQUAL          :  '<='  ;
T_REGEXP             :  '=~'  ;
T_NEQREGEXP          :  '!~'  ;
T_COMMA              :  ','   ;
T_OPEN_B             :  '{'   ;
T_CLOSE_B            :  '}'   ;
T_OPEN_SB            :  '['   ;
T_CLOSE_SB           :  ']'   ;
T_OPEN_P             :  '('   ;
T_CLOSE_P            :  ')'   ;
T_ADD                :  '+'   ;
T_SUB                :  '-'   ;
T_DIV                :  '/'   ;
T_MUL                :  '*'   ;
T_MOD                :  '%'   ;
T_UNDERLINE          :  '_'   ;

L_ID                 : L_ID_PART ;
L_INT                : L_DIGIT+;                                               // Integer
L_DEC                : L_DIGIT+ '.' ~'.' L_DIGIT*                               // Decimal number
                     | '.' L_DIGIT+
                     ;

fragment BLANK       : [ \t\r\n]      ;
fragment L_DIGIT     : [0-9] ;
// Double quoted string escape sequence
//fragment L_STR_ESC_D : '""' | '\\"' ;
fragment L_ID_PART   :
                      [a-zA-Z] ([a-zA-Z] | L_DIGIT | '_' | '.')*                                            // Identifier part
                      | '$' '{' .*? '}'
                      | ('_' | '@' | '#' | '$') ([a-zA-Z] | L_DIGIT | '_' | '@' | ':' | '#' | '$')+     // (at least one char must follow special char)
                      | '"' .*? '"'                                                                           // Quoted identifiers
                      | '`' .*? '`'                                                                           // Quoted identifiers
                      | '\'' .*? '\''                                                                           // Quoted identifiers
                     ;

// Support case-insensitive keywords and allowing case-sensitive identifiers
fragment A : ('a'|'A') ;
fragment B : ('b'|'B') ;
fragment C : ('c'|'C') ;
fragment D : ('d'|'D') ;
fragment E : ('e'|'E') ;
fragment F : ('f'|'F') ;
fragment G : ('g'|'G') ;
fragment H : ('h'|'H') ;
fragment I : ('i'|'I') ;
fragment J : ('j'|'J') ;
fragment K : ('k'|'K') ;
fragment L : ('l'|'L') ;
fragment M : ('m'|'M') ;
fragment N : ('n'|'N') ;
fragment O : ('o'|'O') ;
fragment P : ('p'|'P') ;
fragment Q : ('q'|'Q') ;
fragment R : ('r'|'R') ;
fragment S : ('s'|'S') ;
fragment T : ('t'|'T') ;
fragment U : ('u'|'U') ;
fragment V : ('v'|'V') ;
fragment W : ('w'|'W') ;
fragment X : ('x'|'X') ;
fragment Y : ('y'|'Y') ;
fragment Z : ('z'|'Z') ;
