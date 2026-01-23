import { useQuery } from '@tanstack/react-query';
import { getFeatures } from '../lib/api';

export function useFeatures() {
  return useQuery({
    queryKey: ['features'],
    queryFn: getFeatures,
    refetchInterval: 5000, // Poll every 5 seconds
  });
}
