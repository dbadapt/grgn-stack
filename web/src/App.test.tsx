import { MantineProvider } from '@mantine/core';
import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import App from './App';

describe('App', () => {
  it('renders the app title', () => {
    render(
      <MantineProvider>
        <App />
      </MantineProvider>,
    );

    expect(screen.getByText('GRGN Stack')).toBeInTheDocument();
  });

  it('renders the Google login button', () => {
    render(
      <MantineProvider>
        <App />
      </MantineProvider>,
    );

    expect(screen.getByText('Login with Google')).toBeInTheDocument();
  });

  it('renders the get started button', () => {
    render(
      <MantineProvider>
        <App />
      </MantineProvider>,
    );

    expect(screen.getByText('Get Started')).toBeInTheDocument();
  });
});
