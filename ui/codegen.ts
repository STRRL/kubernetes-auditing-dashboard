import { CodegenConfig } from '@graphql-codegen/cli'

const config: CodegenConfig = {
    schema: [
        "../gql/*.graphql"
    ],
    documents: [
        'components/**/*.tsx',
        'pages/**/*.tsx',
    ],
    ignoreNoDocuments: false, // for better experience with the watcher
    generates: {
        './modules/gql/': {
            preset: 'client',
            plugins: []
        }
    }
}

export default config
