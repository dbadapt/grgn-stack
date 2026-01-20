import { Button, Container, Group, Stack, Text, Title } from '@mantine/core';
import { IconBrandReact } from '@tabler/icons-react';

function App() {
  return (
    <Container size="md" mt="xl">
      <Stack align="center" gap="md">
        <IconBrandReact size={64} stroke={1.5} />
        <Title order={1}>GRGN Stack</Title>
        <Text c="dimmed" size="sm" ta="center" fs="italic">
          (pronounced "Gur-gen")
        </Text>
        <Text c="dimmed" size="lg" ta="center" mt="xs">
          <strong>G</strong>o + <strong>R</strong>eact + <strong>G</strong>
          raphQL + <strong>N</strong>eo4j
        </Text>
        <Text c="dimmed" ta="center">
          A modern, production-ready full-stack template for building
          applications with Go, Neo4j graph database, GraphQL API, and React
          frontend.
        </Text>

        <Group>
          <Button variant="filled" color="blue">
            Login with Google
          </Button>
          <Button variant="light">Get Started</Button>
        </Group>
      </Stack>
    </Container>
  );
}

export default App;
