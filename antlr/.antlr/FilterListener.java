// Generated from /home/jgautier/git-data/poc-fiber/antlr/Filter.g4 by ANTLR 4.13.1


import org.antlr.v4.runtime.tree.ParseTreeListener;

/**
 * This interface defines a complete listener for a parse tree produced by
 * {@link FilterParser}.
 */
public interface FilterListener extends ParseTreeListener {
	/**
	 * Enter a parse tree produced by {@link FilterParser#filter}.
	 * @param ctx the parse tree
	 */
	void enterFilter(FilterParser.FilterContext ctx);
	/**
	 * Exit a parse tree produced by {@link FilterParser#filter}.
	 * @param ctx the parse tree
	 */
	void exitFilter(FilterParser.FilterContext ctx);
	/**
	 * Enter a parse tree produced by {@link FilterParser#expr}.
	 * @param ctx the parse tree
	 */
	void enterExpr(FilterParser.ExprContext ctx);
	/**
	 * Exit a parse tree produced by {@link FilterParser#expr}.
	 * @param ctx the parse tree
	 */
	void exitExpr(FilterParser.ExprContext ctx);
}