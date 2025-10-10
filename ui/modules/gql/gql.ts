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
    "\n  query completedRequestResponseAuditEvents($page: Int, $pageSize: Int){\n    completedRequestResponseAuditEvents(page: $page, pageSize: $pageSize) {\n    total\n    page\n    pageSize\n    totalPages\n    hasNextPage\n    hasPreviousPage\n    rows {\n      id\n      level\n      stage\n      verb\n      useragent\n      resource\n      namespace\n      name\n      stagetimestamp\n      apigroup\n      apiversion\n    }\n  }\n  }\n": types.CompletedRequestResponseAuditEventsDocument,
    "\n  query eventsCount{\n    auditEvents{\n      totalCount\n      pageInfo{\n        hasNextPage\n        hasPreviousPage\n        startCursor\n        endCursor\n      }\n    }\n  }\n": types.EventsCountDocument,
    "\n  query eventsCountNonGet{\n    auditEvents(where: {\n        verb: \"watch\"\n      })\n    {\n      totalCount\n    }\n  }\n": types.EventsCountNonGetDocument,
    "\n  query GetResourceLifecycle(\n    $apiGroup: String!\n    $version: String!\n    $kind: String!\n    $namespace: String\n    $name: String!\n  ) {\n    resourceLifecycle(\n      apiGroup: $apiGroup\n      version: $version\n      kind: $kind\n      namespace: $namespace\n      name: $name\n    ) {\n      id\n      type\n      timestamp\n      user\n      resourceState\n      diff {\n        added\n        removed\n        modified {\n          path\n          oldValue\n          newValue\n        }\n      }\n    }\n  }\n": types.GetResourceLifecycleDocument,
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
export function graphql(source: "\n  query completedRequestResponseAuditEvents($page: Int, $pageSize: Int){\n    completedRequestResponseAuditEvents(page: $page, pageSize: $pageSize) {\n    total\n    page\n    pageSize\n    totalPages\n    hasNextPage\n    hasPreviousPage\n    rows {\n      id\n      level\n      stage\n      verb\n      useragent\n      resource\n      namespace\n      name\n      stagetimestamp\n      apigroup\n      apiversion\n    }\n  }\n  }\n"): (typeof documents)["\n  query completedRequestResponseAuditEvents($page: Int, $pageSize: Int){\n    completedRequestResponseAuditEvents(page: $page, pageSize: $pageSize) {\n    total\n    page\n    pageSize\n    totalPages\n    hasNextPage\n    hasPreviousPage\n    rows {\n      id\n      level\n      stage\n      verb\n      useragent\n      resource\n      namespace\n      name\n      stagetimestamp\n      apigroup\n      apiversion\n    }\n  }\n  }\n"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  query eventsCount{\n    auditEvents{\n      totalCount\n      pageInfo{\n        hasNextPage\n        hasPreviousPage\n        startCursor\n        endCursor\n      }\n    }\n  }\n"): (typeof documents)["\n  query eventsCount{\n    auditEvents{\n      totalCount\n      pageInfo{\n        hasNextPage\n        hasPreviousPage\n        startCursor\n        endCursor\n      }\n    }\n  }\n"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  query eventsCountNonGet{\n    auditEvents(where: {\n        verb: \"watch\"\n      })\n    {\n      totalCount\n    }\n  }\n"): (typeof documents)["\n  query eventsCountNonGet{\n    auditEvents(where: {\n        verb: \"watch\"\n      })\n    {\n      totalCount\n    }\n  }\n"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  query GetResourceLifecycle(\n    $apiGroup: String!\n    $version: String!\n    $kind: String!\n    $namespace: String\n    $name: String!\n  ) {\n    resourceLifecycle(\n      apiGroup: $apiGroup\n      version: $version\n      kind: $kind\n      namespace: $namespace\n      name: $name\n    ) {\n      id\n      type\n      timestamp\n      user\n      resourceState\n      diff {\n        added\n        removed\n        modified {\n          path\n          oldValue\n          newValue\n        }\n      }\n    }\n  }\n"): (typeof documents)["\n  query GetResourceLifecycle(\n    $apiGroup: String!\n    $version: String!\n    $kind: String!\n    $namespace: String\n    $name: String!\n  ) {\n    resourceLifecycle(\n      apiGroup: $apiGroup\n      version: $version\n      kind: $kind\n      namespace: $namespace\n      name: $name\n    ) {\n      id\n      type\n      timestamp\n      user\n      resourceState\n      diff {\n        added\n        removed\n        modified {\n          path\n          oldValue\n          newValue\n        }\n      }\n    }\n  }\n"];

export function graphql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;