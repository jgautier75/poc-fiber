// Generated from /home/jgautier/git-data/poc-fiber/antlr/Filter.g4 by ANTLR 4.13.1


import org.antlr.v4.runtime.Lexer;
import org.antlr.v4.runtime.CharStream;
import org.antlr.v4.runtime.Token;
import org.antlr.v4.runtime.TokenStream;
import org.antlr.v4.runtime.*;
import org.antlr.v4.runtime.atn.*;
import org.antlr.v4.runtime.dfa.DFA;
import org.antlr.v4.runtime.misc.*;

@SuppressWarnings({"all", "warnings", "unchecked", "unused", "cast", "CheckReturnValue", "this-escape"})
public class FilterLexer extends Lexer {
	static { RuntimeMetaData.checkVersion("4.13.1", RuntimeMetaData.VERSION); }

	protected static final DFA[] _decisionToDFA;
	protected static final PredictionContextCache _sharedContextCache =
		new PredictionContextCache();
	public static final int
		OPAR=1, CPAR=2, OR=3, AND=4, NOT=5, COMPARISON=6, GT=7, GE=8, LT=9, LE=10, 
		EQ=11, NE=12, LK=13, VALUE=14, TRUE=15, FALSE=16, PROPERTY=17, STRING=18, 
		INT=19, FLOAT=20, SPACE=21;
	public static String[] channelNames = {
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN"
	};

	public static String[] modeNames = {
		"DEFAULT_MODE"
	};

	private static String[] makeRuleNames() {
		return new String[] {
			"OPAR", "CPAR", "OR", "AND", "NOT", "COMPARISON", "GT", "GE", "LT", "LE", 
			"EQ", "NE", "LK", "VALUE", "TRUE", "FALSE", "PROPERTY", "ALLOWED_CHARACTERS", 
			"STRING", "INT", "FLOAT", "DIGIT", "SPACE"
		};
	}
	public static final String[] ruleNames = makeRuleNames();

	private static String[] makeLiteralNames() {
		return new String[] {
			null, "'('", "')'", "'or'", "'and'", "'not'", null, "'gt'", "'ge'", "'lt'", 
			"'le'", "'eq'", "'ne'", "'lk'", null, "'true'", "'false'"
		};
	}
	private static final String[] _LITERAL_NAMES = makeLiteralNames();
	private static String[] makeSymbolicNames() {
		return new String[] {
			null, "OPAR", "CPAR", "OR", "AND", "NOT", "COMPARISON", "GT", "GE", "LT", 
			"LE", "EQ", "NE", "LK", "VALUE", "TRUE", "FALSE", "PROPERTY", "STRING", 
			"INT", "FLOAT", "SPACE"
		};
	}
	private static final String[] _SYMBOLIC_NAMES = makeSymbolicNames();
	public static final Vocabulary VOCABULARY = new VocabularyImpl(_LITERAL_NAMES, _SYMBOLIC_NAMES);

	/**
	 * @deprecated Use {@link #VOCABULARY} instead.
	 */
	@Deprecated
	public static final String[] tokenNames;
	static {
		tokenNames = new String[_SYMBOLIC_NAMES.length];
		for (int i = 0; i < tokenNames.length; i++) {
			tokenNames[i] = VOCABULARY.getLiteralName(i);
			if (tokenNames[i] == null) {
				tokenNames[i] = VOCABULARY.getSymbolicName(i);
			}

			if (tokenNames[i] == null) {
				tokenNames[i] = "<INVALID>";
			}
		}
	}

	@Override
	@Deprecated
	public String[] getTokenNames() {
		return tokenNames;
	}

	@Override

	public Vocabulary getVocabulary() {
		return VOCABULARY;
	}


	public FilterLexer(CharStream input) {
		super(input);
		_interp = new LexerATNSimulator(this,_ATN,_decisionToDFA,_sharedContextCache);
	}

	@Override
	public String getGrammarFileName() { return "Filter.g4"; }

	@Override
	public String[] getRuleNames() { return ruleNames; }

	@Override
	public String getSerializedATN() { return _serializedATN; }

	@Override
	public String[] getChannelNames() { return channelNames; }

	@Override
	public String[] getModeNames() { return modeNames; }

	@Override
	public ATN getATN() { return _ATN; }

	public static final String _serializedATN =
		"\u0004\u0000\u0015\u00aa\u0006\uffff\uffff\u0002\u0000\u0007\u0000\u0002"+
		"\u0001\u0007\u0001\u0002\u0002\u0007\u0002\u0002\u0003\u0007\u0003\u0002"+
		"\u0004\u0007\u0004\u0002\u0005\u0007\u0005\u0002\u0006\u0007\u0006\u0002"+
		"\u0007\u0007\u0007\u0002\b\u0007\b\u0002\t\u0007\t\u0002\n\u0007\n\u0002"+
		"\u000b\u0007\u000b\u0002\f\u0007\f\u0002\r\u0007\r\u0002\u000e\u0007\u000e"+
		"\u0002\u000f\u0007\u000f\u0002\u0010\u0007\u0010\u0002\u0011\u0007\u0011"+
		"\u0002\u0012\u0007\u0012\u0002\u0013\u0007\u0013\u0002\u0014\u0007\u0014"+
		"\u0002\u0015\u0007\u0015\u0002\u0016\u0007\u0016\u0001\u0000\u0001\u0000"+
		"\u0001\u0001\u0001\u0001\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0003"+
		"\u0001\u0003\u0001\u0003\u0001\u0003\u0001\u0004\u0001\u0004\u0001\u0004"+
		"\u0001\u0004\u0001\u0005\u0001\u0005\u0001\u0005\u0001\u0005\u0001\u0005"+
		"\u0001\u0005\u0001\u0005\u0003\u0005F\b\u0005\u0001\u0006\u0001\u0006"+
		"\u0001\u0006\u0001\u0007\u0001\u0007\u0001\u0007\u0001\b\u0001\b\u0001"+
		"\b\u0001\t\u0001\t\u0001\t\u0001\n\u0001\n\u0001\n\u0001\u000b\u0001\u000b"+
		"\u0001\u000b\u0001\f\u0001\f\u0001\f\u0001\r\u0001\r\u0001\r\u0001\r\u0001"+
		"\r\u0003\rb\b\r\u0001\u000e\u0001\u000e\u0001\u000e\u0001\u000e\u0001"+
		"\u000e\u0001\u000f\u0001\u000f\u0001\u000f\u0001\u000f\u0001\u000f\u0001"+
		"\u000f\u0001\u0010\u0004\u0010p\b\u0010\u000b\u0010\f\u0010q\u0001\u0010"+
		"\u0001\u0010\u0004\u0010v\b\u0010\u000b\u0010\f\u0010w\u0005\u0010z\b"+
		"\u0010\n\u0010\f\u0010}\t\u0010\u0001\u0011\u0001\u0011\u0001\u0012\u0001"+
		"\u0012\u0001\u0012\u0001\u0012\u0005\u0012\u0085\b\u0012\n\u0012\f\u0012"+
		"\u0088\t\u0012\u0001\u0012\u0001\u0012\u0001\u0013\u0004\u0013\u008d\b"+
		"\u0013\u000b\u0013\f\u0013\u008e\u0001\u0014\u0004\u0014\u0092\b\u0014"+
		"\u000b\u0014\f\u0014\u0093\u0001\u0014\u0001\u0014\u0005\u0014\u0098\b"+
		"\u0014\n\u0014\f\u0014\u009b\t\u0014\u0001\u0014\u0001\u0014\u0004\u0014"+
		"\u009f\b\u0014\u000b\u0014\f\u0014\u00a0\u0003\u0014\u00a3\b\u0014\u0001"+
		"\u0015\u0001\u0015\u0001\u0016\u0001\u0016\u0001\u0016\u0001\u0016\u0000"+
		"\u0000\u0017\u0001\u0001\u0003\u0002\u0005\u0003\u0007\u0004\t\u0005\u000b"+
		"\u0006\r\u0007\u000f\b\u0011\t\u0013\n\u0015\u000b\u0017\f\u0019\r\u001b"+
		"\u000e\u001d\u000f\u001f\u0010!\u0011#\u0000%\u0012\'\u0013)\u0014+\u0000"+
		"-\u0015\u0001\u0000\u0004\u0005\u0000--09AZ__az\u0003\u0000\n\n\r\r\'"+
		"\'\u0001\u000009\u0003\u0000\t\n\r\r  \u00bb\u0000\u0001\u0001\u0000\u0000"+
		"\u0000\u0000\u0003\u0001\u0000\u0000\u0000\u0000\u0005\u0001\u0000\u0000"+
		"\u0000\u0000\u0007\u0001\u0000\u0000\u0000\u0000\t\u0001\u0000\u0000\u0000"+
		"\u0000\u000b\u0001\u0000\u0000\u0000\u0000\r\u0001\u0000\u0000\u0000\u0000"+
		"\u000f\u0001\u0000\u0000\u0000\u0000\u0011\u0001\u0000\u0000\u0000\u0000"+
		"\u0013\u0001\u0000\u0000\u0000\u0000\u0015\u0001\u0000\u0000\u0000\u0000"+
		"\u0017\u0001\u0000\u0000\u0000\u0000\u0019\u0001\u0000\u0000\u0000\u0000"+
		"\u001b\u0001\u0000\u0000\u0000\u0000\u001d\u0001\u0000\u0000\u0000\u0000"+
		"\u001f\u0001\u0000\u0000\u0000\u0000!\u0001\u0000\u0000\u0000\u0000%\u0001"+
		"\u0000\u0000\u0000\u0000\'\u0001\u0000\u0000\u0000\u0000)\u0001\u0000"+
		"\u0000\u0000\u0000-\u0001\u0000\u0000\u0000\u0001/\u0001\u0000\u0000\u0000"+
		"\u00031\u0001\u0000\u0000\u0000\u00053\u0001\u0000\u0000\u0000\u00076"+
		"\u0001\u0000\u0000\u0000\t:\u0001\u0000\u0000\u0000\u000bE\u0001\u0000"+
		"\u0000\u0000\rG\u0001\u0000\u0000\u0000\u000fJ\u0001\u0000\u0000\u0000"+
		"\u0011M\u0001\u0000\u0000\u0000\u0013P\u0001\u0000\u0000\u0000\u0015S"+
		"\u0001\u0000\u0000\u0000\u0017V\u0001\u0000\u0000\u0000\u0019Y\u0001\u0000"+
		"\u0000\u0000\u001ba\u0001\u0000\u0000\u0000\u001dc\u0001\u0000\u0000\u0000"+
		"\u001fh\u0001\u0000\u0000\u0000!o\u0001\u0000\u0000\u0000#~\u0001\u0000"+
		"\u0000\u0000%\u0080\u0001\u0000\u0000\u0000\'\u008c\u0001\u0000\u0000"+
		"\u0000)\u00a2\u0001\u0000\u0000\u0000+\u00a4\u0001\u0000\u0000\u0000-"+
		"\u00a6\u0001\u0000\u0000\u0000/0\u0005(\u0000\u00000\u0002\u0001\u0000"+
		"\u0000\u000012\u0005)\u0000\u00002\u0004\u0001\u0000\u0000\u000034\u0005"+
		"o\u0000\u000045\u0005r\u0000\u00005\u0006\u0001\u0000\u0000\u000067\u0005"+
		"a\u0000\u000078\u0005n\u0000\u000089\u0005d\u0000\u00009\b\u0001\u0000"+
		"\u0000\u0000:;\u0005n\u0000\u0000;<\u0005o\u0000\u0000<=\u0005t\u0000"+
		"\u0000=\n\u0001\u0000\u0000\u0000>F\u0003\r\u0006\u0000?F\u0003\u000f"+
		"\u0007\u0000@F\u0003\u0011\b\u0000AF\u0003\u0013\t\u0000BF\u0003\u0015"+
		"\n\u0000CF\u0003\u0017\u000b\u0000DF\u0003\u0019\f\u0000E>\u0001\u0000"+
		"\u0000\u0000E?\u0001\u0000\u0000\u0000E@\u0001\u0000\u0000\u0000EA\u0001"+
		"\u0000\u0000\u0000EB\u0001\u0000\u0000\u0000EC\u0001\u0000\u0000\u0000"+
		"ED\u0001\u0000\u0000\u0000F\f\u0001\u0000\u0000\u0000GH\u0005g\u0000\u0000"+
		"HI\u0005t\u0000\u0000I\u000e\u0001\u0000\u0000\u0000JK\u0005g\u0000\u0000"+
		"KL\u0005e\u0000\u0000L\u0010\u0001\u0000\u0000\u0000MN\u0005l\u0000\u0000"+
		"NO\u0005t\u0000\u0000O\u0012\u0001\u0000\u0000\u0000PQ\u0005l\u0000\u0000"+
		"QR\u0005e\u0000\u0000R\u0014\u0001\u0000\u0000\u0000ST\u0005e\u0000\u0000"+
		"TU\u0005q\u0000\u0000U\u0016\u0001\u0000\u0000\u0000VW\u0005n\u0000\u0000"+
		"WX\u0005e\u0000\u0000X\u0018\u0001\u0000\u0000\u0000YZ\u0005l\u0000\u0000"+
		"Z[\u0005k\u0000\u0000[\u001a\u0001\u0000\u0000\u0000\\b\u0003\u001d\u000e"+
		"\u0000]b\u0003\u001f\u000f\u0000^b\u0003\'\u0013\u0000_b\u0003)\u0014"+
		"\u0000`b\u0003%\u0012\u0000a\\\u0001\u0000\u0000\u0000a]\u0001\u0000\u0000"+
		"\u0000a^\u0001\u0000\u0000\u0000a_\u0001\u0000\u0000\u0000a`\u0001\u0000"+
		"\u0000\u0000b\u001c\u0001\u0000\u0000\u0000cd\u0005t\u0000\u0000de\u0005"+
		"r\u0000\u0000ef\u0005u\u0000\u0000fg\u0005e\u0000\u0000g\u001e\u0001\u0000"+
		"\u0000\u0000hi\u0005f\u0000\u0000ij\u0005a\u0000\u0000jk\u0005l\u0000"+
		"\u0000kl\u0005s\u0000\u0000lm\u0005e\u0000\u0000m \u0001\u0000\u0000\u0000"+
		"np\u0003#\u0011\u0000on\u0001\u0000\u0000\u0000pq\u0001\u0000\u0000\u0000"+
		"qo\u0001\u0000\u0000\u0000qr\u0001\u0000\u0000\u0000r{\u0001\u0000\u0000"+
		"\u0000su\u0005.\u0000\u0000tv\u0003#\u0011\u0000ut\u0001\u0000\u0000\u0000"+
		"vw\u0001\u0000\u0000\u0000wu\u0001\u0000\u0000\u0000wx\u0001\u0000\u0000"+
		"\u0000xz\u0001\u0000\u0000\u0000ys\u0001\u0000\u0000\u0000z}\u0001\u0000"+
		"\u0000\u0000{y\u0001\u0000\u0000\u0000{|\u0001\u0000\u0000\u0000|\"\u0001"+
		"\u0000\u0000\u0000}{\u0001\u0000\u0000\u0000~\u007f\u0007\u0000\u0000"+
		"\u0000\u007f$\u0001\u0000\u0000\u0000\u0080\u0086\u0005\'\u0000\u0000"+
		"\u0081\u0085\b\u0001\u0000\u0000\u0082\u0083\u0005\'\u0000\u0000\u0083"+
		"\u0085\u0005\'\u0000\u0000\u0084\u0081\u0001\u0000\u0000\u0000\u0084\u0082"+
		"\u0001\u0000\u0000\u0000\u0085\u0088\u0001\u0000\u0000\u0000\u0086\u0084"+
		"\u0001\u0000\u0000\u0000\u0086\u0087\u0001\u0000\u0000\u0000\u0087\u0089"+
		"\u0001\u0000\u0000\u0000\u0088\u0086\u0001\u0000\u0000\u0000\u0089\u008a"+
		"\u0005\'\u0000\u0000\u008a&\u0001\u0000\u0000\u0000\u008b\u008d\u0003"+
		"+\u0015\u0000\u008c\u008b\u0001\u0000\u0000\u0000\u008d\u008e\u0001\u0000"+
		"\u0000\u0000\u008e\u008c\u0001\u0000\u0000\u0000\u008e\u008f\u0001\u0000"+
		"\u0000\u0000\u008f(\u0001\u0000\u0000\u0000\u0090\u0092\u0003+\u0015\u0000"+
		"\u0091\u0090\u0001\u0000\u0000\u0000\u0092\u0093\u0001\u0000\u0000\u0000"+
		"\u0093\u0091\u0001\u0000\u0000\u0000\u0093\u0094\u0001\u0000\u0000\u0000"+
		"\u0094\u0095\u0001\u0000\u0000\u0000\u0095\u0099\u0005.\u0000\u0000\u0096"+
		"\u0098\u0003+\u0015\u0000\u0097\u0096\u0001\u0000\u0000\u0000\u0098\u009b"+
		"\u0001\u0000\u0000\u0000\u0099\u0097\u0001\u0000\u0000\u0000\u0099\u009a"+
		"\u0001\u0000\u0000\u0000\u009a\u00a3\u0001\u0000\u0000\u0000\u009b\u0099"+
		"\u0001\u0000\u0000\u0000\u009c\u009e\u0005.\u0000\u0000\u009d\u009f\u0003"+
		"+\u0015\u0000\u009e\u009d\u0001\u0000\u0000\u0000\u009f\u00a0\u0001\u0000"+
		"\u0000\u0000\u00a0\u009e\u0001\u0000\u0000\u0000\u00a0\u00a1\u0001\u0000"+
		"\u0000\u0000\u00a1\u00a3\u0001\u0000\u0000\u0000\u00a2\u0091\u0001\u0000"+
		"\u0000\u0000\u00a2\u009c\u0001\u0000\u0000\u0000\u00a3*\u0001\u0000\u0000"+
		"\u0000\u00a4\u00a5\u0007\u0002\u0000\u0000\u00a5,\u0001\u0000\u0000\u0000"+
		"\u00a6\u00a7\u0007\u0003\u0000\u0000\u00a7\u00a8\u0001\u0000\u0000\u0000"+
		"\u00a8\u00a9\u0006\u0016\u0000\u0000\u00a9.\u0001\u0000\u0000\u0000\r"+
		"\u0000Eaqw{\u0084\u0086\u008e\u0093\u0099\u00a0\u00a2\u0001\u0006\u0000"+
		"\u0000";
	public static final ATN _ATN =
		new ATNDeserializer().deserialize(_serializedATN.toCharArray());
	static {
		_decisionToDFA = new DFA[_ATN.getNumberOfDecisions()];
		for (int i = 0; i < _ATN.getNumberOfDecisions(); i++) {
			_decisionToDFA[i] = new DFA(_ATN.getDecisionState(i), i);
		}
	}
}