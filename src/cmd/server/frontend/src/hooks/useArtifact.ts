import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getArtifact, saveArtifact } from '../lib/api';

export function useArtifact(featureId: string | null, artifactType: string) {
  return useQuery({
    queryKey: ['artifact', featureId, artifactType],
    queryFn: () => getArtifact(featureId!, artifactType),
    enabled: !!featureId,
  });
}

export function useSaveArtifact(featureId: string | null, artifactType: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (content: string) => saveArtifact(featureId!, artifactType, content),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['artifact', featureId, artifactType] });
      queryClient.invalidateQueries({ queryKey: ['features'] });
    },
  });
}
