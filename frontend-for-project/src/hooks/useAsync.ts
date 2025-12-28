import { useState, useEffect } from 'react';

interface UseAsyncState<T> {
  data: T | null;
  isLoading: boolean;
  error: Error | null;
}

export const useAsync = <T,>(
  asyncFunction: () => Promise<T>,
  immediate = true,
  dependencies: any[] = []
): UseAsyncState<T> & { execute: () => Promise<T> } => {
  const [state, setState] = useState<UseAsyncState<T>>({
    data: null,
    isLoading: immediate,
    error: null,
  });

  const execute = async () => {
    setState({ data: null, isLoading: true, error: null });
    try {
      const response = await asyncFunction();
      setState({ data: response, isLoading: false, error: null });
      return response;
    } catch (error) {
      setState({
        data: null,
        isLoading: false,
        error: error instanceof Error ? error : new Error('Unknown error'),
      });
      throw error;
    }
  };

  useEffect(() => {
    if (immediate) {
      execute();
    }
  }, dependencies);

  return { ...state, execute };
};
