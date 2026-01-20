import {
  useQuery,
  useInfiniteQuery,
  UseQueryOptions,
  UseInfiniteQueryOptions,
} from '@tanstack/react-query';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = {
  [K in keyof T]: T[K];
};
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & {
  [SubKey in K]?: Maybe<T[SubKey]>;
};
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & {
  [SubKey in K]: Maybe<T[SubKey]>;
};
export type MakeEmpty<
  T extends { [key: string]: unknown },
  K extends keyof T,
> = { [_ in K]?: never };
export type Incremental<T> =
  | T
  | {
      [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never;
    };

function fetcher<TData, TVariables>(query: string, variables?: TVariables) {
  return async (): Promise<TData> => {
    const res = await fetch('http://localhost:8080/graphql', {
      method: 'POST',
      ...{ headers: { 'Content-Type': 'application/json' } },
      body: JSON.stringify({ query, variables }),
    });

    const json = await res.json();

    if (json.errors) {
      const { message } = json.errors[0];

      throw new Error(message);
    }

    return json.data;
  };
}
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string };
  String: { input: string; output: string };
  Boolean: { input: boolean; output: boolean };
  Int: { input: number; output: number };
  Float: { input: number; output: number };
  Time: { input: string; output: string };
  UUID: { input: string; output: string };
};

export type Mutation = {
  __typename?: 'Mutation';
  _empty?: Maybe<Scalars['String']['output']>;
};

export type PaginationInput = {
  page?: InputMaybe<Scalars['Int']['input']>;
  pageSize?: InputMaybe<Scalars['Int']['input']>;
};

export type Query = {
  __typename?: 'Query';
  health: Scalars['String']['output'];
  me?: Maybe<User>;
  user?: Maybe<User>;
};

export type QueryUserArgs = {
  id: Scalars['ID']['input'];
};

export type Subscription = {
  __typename?: 'Subscription';
  _empty?: Maybe<Scalars['String']['output']>;
};

export type User = {
  __typename?: 'User';
  createdAt: Scalars['Time']['output'];
  email: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  name?: Maybe<Scalars['String']['output']>;
  updatedAt: Scalars['Time']['output'];
};

export type GetHealthQueryVariables = Exact<{ [key: string]: never }>;

export type GetHealthQuery = { __typename?: 'Query'; health: string };

export type GetCurrentUserQueryVariables = Exact<{ [key: string]: never }>;

export type GetCurrentUserQuery = {
  __typename?: 'Query';
  me?: {
    __typename?: 'User';
    id: string;
    email: string;
    name?: string | null;
    createdAt: string;
    updatedAt: string;
  } | null;
};

export type GetUserQueryVariables = Exact<{
  id: Scalars['ID']['input'];
}>;

export type GetUserQuery = {
  __typename?: 'Query';
  user?: {
    __typename?: 'User';
    id: string;
    email: string;
    name?: string | null;
    createdAt: string;
    updatedAt: string;
  } | null;
};

export const GetHealthDocument = `
    query GetHealth {
  health
}
    `;

export const useGetHealthQuery = <TData = GetHealthQuery, TError = unknown>(
  variables?: GetHealthQueryVariables,
  options?: UseQueryOptions<GetHealthQuery, TError, TData>,
) => {
  return useQuery<GetHealthQuery, TError, TData>(
    variables === undefined ? ['GetHealth'] : ['GetHealth', variables],
    fetcher<GetHealthQuery, GetHealthQueryVariables>(
      GetHealthDocument,
      variables,
    ),
    options,
  );
};

useGetHealthQuery.getKey = (variables?: GetHealthQueryVariables) =>
  variables === undefined ? ['GetHealth'] : ['GetHealth', variables];

export const useInfiniteGetHealthQuery = <
  TData = GetHealthQuery,
  TError = unknown,
>(
  variables?: GetHealthQueryVariables,
  options?: UseInfiniteQueryOptions<GetHealthQuery, TError, TData>,
) => {
  return useInfiniteQuery<GetHealthQuery, TError, TData>(
    variables === undefined
      ? ['GetHealth.infinite']
      : ['GetHealth.infinite', variables],
    (metaData) =>
      fetcher<GetHealthQuery, GetHealthQueryVariables>(GetHealthDocument, {
        ...variables,
        ...(metaData.pageParam ?? {}),
      })(),
    options,
  );
};

useInfiniteGetHealthQuery.getKey = (variables?: GetHealthQueryVariables) =>
  variables === undefined
    ? ['GetHealth.infinite']
    : ['GetHealth.infinite', variables];

useGetHealthQuery.fetcher = (variables?: GetHealthQueryVariables) =>
  fetcher<GetHealthQuery, GetHealthQueryVariables>(
    GetHealthDocument,
    variables,
  );

export const GetCurrentUserDocument = `
    query GetCurrentUser {
  me {
    id
    email
    name
    createdAt
    updatedAt
  }
}
    `;

export const useGetCurrentUserQuery = <
  TData = GetCurrentUserQuery,
  TError = unknown,
>(
  variables?: GetCurrentUserQueryVariables,
  options?: UseQueryOptions<GetCurrentUserQuery, TError, TData>,
) => {
  return useQuery<GetCurrentUserQuery, TError, TData>(
    variables === undefined
      ? ['GetCurrentUser']
      : ['GetCurrentUser', variables],
    fetcher<GetCurrentUserQuery, GetCurrentUserQueryVariables>(
      GetCurrentUserDocument,
      variables,
    ),
    options,
  );
};

useGetCurrentUserQuery.getKey = (variables?: GetCurrentUserQueryVariables) =>
  variables === undefined ? ['GetCurrentUser'] : ['GetCurrentUser', variables];

export const useInfiniteGetCurrentUserQuery = <
  TData = GetCurrentUserQuery,
  TError = unknown,
>(
  variables?: GetCurrentUserQueryVariables,
  options?: UseInfiniteQueryOptions<GetCurrentUserQuery, TError, TData>,
) => {
  return useInfiniteQuery<GetCurrentUserQuery, TError, TData>(
    variables === undefined
      ? ['GetCurrentUser.infinite']
      : ['GetCurrentUser.infinite', variables],
    (metaData) =>
      fetcher<GetCurrentUserQuery, GetCurrentUserQueryVariables>(
        GetCurrentUserDocument,
        { ...variables, ...(metaData.pageParam ?? {}) },
      )(),
    options,
  );
};

useInfiniteGetCurrentUserQuery.getKey = (
  variables?: GetCurrentUserQueryVariables,
) =>
  variables === undefined
    ? ['GetCurrentUser.infinite']
    : ['GetCurrentUser.infinite', variables];

useGetCurrentUserQuery.fetcher = (variables?: GetCurrentUserQueryVariables) =>
  fetcher<GetCurrentUserQuery, GetCurrentUserQueryVariables>(
    GetCurrentUserDocument,
    variables,
  );

export const GetUserDocument = `
    query GetUser($id: ID!) {
  user(id: $id) {
    id
    email
    name
    createdAt
    updatedAt
  }
}
    `;

export const useGetUserQuery = <TData = GetUserQuery, TError = unknown>(
  variables: GetUserQueryVariables,
  options?: UseQueryOptions<GetUserQuery, TError, TData>,
) => {
  return useQuery<GetUserQuery, TError, TData>(
    ['GetUser', variables],
    fetcher<GetUserQuery, GetUserQueryVariables>(GetUserDocument, variables),
    options,
  );
};

useGetUserQuery.getKey = (variables: GetUserQueryVariables) => [
  'GetUser',
  variables,
];

export const useInfiniteGetUserQuery = <TData = GetUserQuery, TError = unknown>(
  variables: GetUserQueryVariables,
  options?: UseInfiniteQueryOptions<GetUserQuery, TError, TData>,
) => {
  return useInfiniteQuery<GetUserQuery, TError, TData>(
    ['GetUser.infinite', variables],
    (metaData) =>
      fetcher<GetUserQuery, GetUserQueryVariables>(GetUserDocument, {
        ...variables,
        ...(metaData.pageParam ?? {}),
      })(),
    options,
  );
};

useInfiniteGetUserQuery.getKey = (variables: GetUserQueryVariables) => [
  'GetUser.infinite',
  variables,
];

useGetUserQuery.fetcher = (variables: GetUserQueryVariables) =>
  fetcher<GetUserQuery, GetUserQueryVariables>(GetUserDocument, variables);
