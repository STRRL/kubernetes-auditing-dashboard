/* eslint-disable */
import * as types from './graphql';
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';

/**
 * Map of all GraphQL operations in the project.
 *
 * This map has several performance disadvantages:
 * 1. It is not tree-shakeable, so it will include all operations in the project.
 * 2. It is not minifiable, so the string of a GraphQL query will be multiple times inside the bundle.
 * 3. It does not support dead code elimination, so it will add unused operations.
 *
 * Therefore it is highly recommended to use the babel or swc plugin for production.
 */
const documents = {
    "\nquery eventsCount{\n  auditEvents{\n    totalCount\n    pageInfo{\n      hasNextPage\n      hasPreviousPage\n      startCursor\n      endCursor\n    }\n  }\n}\n": types.EventsCountDocument,
    "\nquery eventsList($first: Int, $after: Cursor){\n  auditEvents(first: $first, after: $after){\n    totalCount\n    pageInfo{\n      hasNextPage\n      hasPreviousPage\n      startCursor\n      endCursor\n    }\n    edges{\n      node{\n        id\n        level\n        stage\n        verb\n        useragent\n        resource\n        namespace\n        name\n        stagetimestamp\n      }\n    }\n  }}\n": types.EventsListDocument,
};

/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 *
 *
 * @example
 * ```ts
 * const query = graphql(`query GetUser($id: ID!) { user(id: $id) { name } }`);
 * ```
 *
 * The query argument is unknown!
 * Please regenerate the types.
 */
export function graphql(source: string): unknown;

/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\nquery eventsCount{\n  auditEvents{\n    totalCount\n    pageInfo{\n      hasNextPage\n      hasPreviousPage\n      startCursor\n      endCursor\n    }\n  }\n}\n"): (typeof documents)["\nquery eventsCount{\n  auditEvents{\n    totalCount\n    pageInfo{\n      hasNextPage\n      hasPreviousPage\n      startCursor\n      endCursor\n    }\n  }\n}\n"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\nquery eventsList($first: Int, $after: Cursor){\n  auditEvents(first: $first, after: $after){\n    totalCount\n    pageInfo{\n      hasNextPage\n      hasPreviousPage\n      startCursor\n      endCursor\n    }\n    edges{\n      node{\n        id\n        level\n        stage\n        verb\n        useragent\n        resource\n        namespace\n        name\n        stagetimestamp\n      }\n    }\n  }}\n"): (typeof documents)["\nquery eventsList($first: Int, $after: Cursor){\n  auditEvents(first: $first, after: $after){\n    totalCount\n    pageInfo{\n      hasNextPage\n      hasPreviousPage\n      startCursor\n      endCursor\n    }\n    edges{\n      node{\n        id\n        level\n        stage\n        verb\n        useragent\n        resource\n        namespace\n        name\n        stagetimestamp\n      }\n    }\n  }}\n"];

export function graphql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;