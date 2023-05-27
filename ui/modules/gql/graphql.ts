/* eslint-disable */
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
  /**
   * Define a Relay Cursor type:
   * https://relay.dev/graphql/connections.htm#sec-Cursor
   */
  Cursor: any;
  /** The builtin Time type */
  Time: any;
};

export type AuditEvent = Node & {
  __typename?: 'AuditEvent';
  apigroup: Scalars['String'];
  apiversion: Scalars['String'];
  auditid: Scalars['String'];
  id: Scalars['ID'];
  level: Scalars['String'];
  name: Scalars['String'];
  namespace: Scalars['String'];
  raw: Scalars['String'];
  requesttimestamp: Scalars['Time'];
  resource: Scalars['String'];
  stage: Scalars['String'];
  stagetimestamp: Scalars['Time'];
  subresource: Scalars['String'];
  useragent: Scalars['String'];
  verb: Scalars['String'];
};

/** A connection to a list of items. */
export type AuditEventConnection = {
  __typename?: 'AuditEventConnection';
  /** A list of edges. */
  edges?: Maybe<Array<Maybe<AuditEventEdge>>>;
  /** Information to aid in pagination. */
  pageInfo: PageInfo;
  /** Identifies the total count of items in the connection. */
  totalCount: Scalars['Int'];
};

/** An edge in a connection. */
export type AuditEventEdge = {
  __typename?: 'AuditEventEdge';
  /** A cursor for use in pagination. */
  cursor: Scalars['Cursor'];
  /** The item at the end of the edge. */
  node?: Maybe<AuditEvent>;
};

export type AuditEventPagination = {
  __typename?: 'AuditEventPagination';
  hasNextPage: Scalars['Boolean'];
  hasPreviousPage: Scalars['Boolean'];
  page: Scalars['Int'];
  pageSize: Scalars['Int'];
  rows: Array<Maybe<AuditEvent>>;
  total: Scalars['Int'];
  totalPages: Scalars['Int'];
};

export type Mutation = {
  __typename?: 'Mutation';
  importResourceKindTSV: Scalars['Int'];
};


export type MutationImportResourceKindTsvArgs = {
  tsv: Scalars['String'];
};

/**
 * An object with an ID.
 * Follows the [Relay Global Object Identification Specification](https://relay.dev/graphql/objectidentification.htm)
 */
export type Node = {
  /** The id of the object. */
  id: Scalars['ID'];
};

/** Possible directions in which to order a list of items when provided an `orderBy` argument. */
export enum OrderDirection {
  /** Specifies an ascending order for a given `orderBy` argument. */
  Asc = 'ASC',
  /** Specifies a descending order for a given `orderBy` argument. */
  Desc = 'DESC'
}

/**
 * Information about pagination in a connection.
 * https://relay.dev/graphql/connections.htm#sec-undefined.PageInfo
 */
export type PageInfo = {
  __typename?: 'PageInfo';
  /** When paginating forwards, the cursor to continue. */
  endCursor?: Maybe<Scalars['Cursor']>;
  /** When paginating forwards, are there more items? */
  hasNextPage: Scalars['Boolean'];
  /** When paginating backwards, are there more items? */
  hasPreviousPage: Scalars['Boolean'];
  /** When paginating backwards, the cursor to continue. */
  startCursor?: Maybe<Scalars['Cursor']>;
};

export type Query = {
  __typename?: 'Query';
  auditEvents: AuditEventConnection;
  completedRequestResponseAuditEvents: AuditEventPagination;
  /** Fetches an object given its ID. */
  node?: Maybe<Node>;
  /** Lookup nodes by a list of IDs. */
  nodes: Array<Maybe<Node>>;
  resourceKinds: ResourceKindConnection;
};


export type QueryAuditEventsArgs = {
  after?: InputMaybe<Scalars['Cursor']>;
  before?: InputMaybe<Scalars['Cursor']>;
  first?: InputMaybe<Scalars['Int']>;
  last?: InputMaybe<Scalars['Int']>;
};


export type QueryCompletedRequestResponseAuditEventsArgs = {
  page?: InputMaybe<Scalars['Int']>;
  pageSize?: InputMaybe<Scalars['Int']>;
};


export type QueryNodeArgs = {
  id: Scalars['ID'];
};


export type QueryNodesArgs = {
  ids: Array<Scalars['ID']>;
};


export type QueryResourceKindsArgs = {
  after?: InputMaybe<Scalars['Cursor']>;
  before?: InputMaybe<Scalars['Cursor']>;
  first?: InputMaybe<Scalars['Int']>;
  last?: InputMaybe<Scalars['Int']>;
};

export type ResourceKind = Node & {
  __typename?: 'ResourceKind';
  apiversion: Scalars['String'];
  id: Scalars['ID'];
  kind: Scalars['String'];
  name: Scalars['String'];
  namespaced: Scalars['Boolean'];
};

/** A connection to a list of items. */
export type ResourceKindConnection = {
  __typename?: 'ResourceKindConnection';
  /** A list of edges. */
  edges?: Maybe<Array<Maybe<ResourceKindEdge>>>;
  /** Information to aid in pagination. */
  pageInfo: PageInfo;
  /** Identifies the total count of items in the connection. */
  totalCount: Scalars['Int'];
};

/** An edge in a connection. */
export type ResourceKindEdge = {
  __typename?: 'ResourceKindEdge';
  /** A cursor for use in pagination. */
  cursor: Scalars['Cursor'];
  /** The item at the end of the edge. */
  node?: Maybe<ResourceKind>;
};

export type View = Node & {
  __typename?: 'View';
  id: Scalars['ID'];
};

export type EventsCountQueryVariables = Exact<{ [key: string]: never; }>;


export type EventsCountQuery = { __typename?: 'Query', auditEvents: { __typename?: 'AuditEventConnection', totalCount: number, pageInfo: { __typename?: 'PageInfo', hasNextPage: boolean, hasPreviousPage: boolean, startCursor?: any | null, endCursor?: any | null } } };

export type CompletedRequestResponseAuditEventsQueryVariables = Exact<{
  page?: InputMaybe<Scalars['Int']>;
  pageSize?: InputMaybe<Scalars['Int']>;
}>;


export type CompletedRequestResponseAuditEventsQuery = { __typename?: 'Query', completedRequestResponseAuditEvents: { __typename?: 'AuditEventPagination', total: number, page: number, pageSize: number, totalPages: number, hasNextPage: boolean, hasPreviousPage: boolean, rows: Array<{ __typename?: 'AuditEvent', id: string, level: string, stage: string, verb: string, useragent: string, resource: string, namespace: string, name: string, stagetimestamp: any, apigroup: string, apiversion: string } | null> } };


export const EventsCountDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"eventsCount"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"auditEvents"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"totalCount"}},{"kind":"Field","name":{"kind":"Name","value":"pageInfo"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"hasNextPage"}},{"kind":"Field","name":{"kind":"Name","value":"hasPreviousPage"}},{"kind":"Field","name":{"kind":"Name","value":"startCursor"}},{"kind":"Field","name":{"kind":"Name","value":"endCursor"}}]}}]}}]}}]} as unknown as DocumentNode<EventsCountQuery, EventsCountQueryVariables>;
export const CompletedRequestResponseAuditEventsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"completedRequestResponseAuditEvents"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"page"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"pageSize"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"completedRequestResponseAuditEvents"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"page"},"value":{"kind":"Variable","name":{"kind":"Name","value":"page"}}},{"kind":"Argument","name":{"kind":"Name","value":"pageSize"},"value":{"kind":"Variable","name":{"kind":"Name","value":"pageSize"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"page"}},{"kind":"Field","name":{"kind":"Name","value":"pageSize"}},{"kind":"Field","name":{"kind":"Name","value":"totalPages"}},{"kind":"Field","name":{"kind":"Name","value":"hasNextPage"}},{"kind":"Field","name":{"kind":"Name","value":"hasPreviousPage"}},{"kind":"Field","name":{"kind":"Name","value":"rows"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"level"}},{"kind":"Field","name":{"kind":"Name","value":"stage"}},{"kind":"Field","name":{"kind":"Name","value":"verb"}},{"kind":"Field","name":{"kind":"Name","value":"useragent"}},{"kind":"Field","name":{"kind":"Name","value":"resource"}},{"kind":"Field","name":{"kind":"Name","value":"namespace"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"stagetimestamp"}},{"kind":"Field","name":{"kind":"Name","value":"apigroup"}},{"kind":"Field","name":{"kind":"Name","value":"apiversion"}}]}}]}}]}}]} as unknown as DocumentNode<CompletedRequestResponseAuditEventsQuery, CompletedRequestResponseAuditEventsQueryVariables>;