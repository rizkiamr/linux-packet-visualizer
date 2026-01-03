import { useState, useEffect, useMemo } from 'react';

/**
 * Hook to load and parse the packet paths contract JSON.
 * Supports multiple paths (egress and ingress).
 * @returns {{ paths, selectedPath, setSelectedPath, simulation, metadata, loading, error }}
 */
export function usePacketPath() {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [selectedPathId, setSelectedPathId] = useState(null);

  useEffect(() => {
    fetch('/data/egress_path.json')
      .then((res) => {
        if (!res.ok) {
          throw new Error(`Failed to load contract: ${res.status}`);
        }
        return res.json();
      })
      .then((json) => {
        setData(json);
        // Default to first path (egress)
        if (json.paths?.length > 0) {
          setSelectedPathId(json.paths[0].path.id);
        }
        setLoading(false);
      })
      .catch((err) => {
        setError(err.message);
        setLoading(false);
      });
  }, []);

  // Get all available paths
  const paths = useMemo(() => {
    return data?.paths?.map(p => ({
      id: p.path.id,
      name: p.path.name,
      direction: p.path.direction,
    })) ?? [];
  }, [data]);

  // Get selected path data
  const selectedPathData = useMemo(() => {
    if (!data?.paths || !selectedPathId) return null;
    return data.paths.find(p => p.path.id === selectedPathId);
  }, [data, selectedPathId]);

  // Set selected path by ID
  const setSelectedPath = (pathId) => {
    setSelectedPathId(pathId);
  };

  return {
    paths,
    selectedPathId,
    setSelectedPath,
    path: selectedPathData?.path ?? null,
    simulation: selectedPathData?.simulation ?? [],
    metadata: data?.metadata ?? null,
    version: data?.version ?? null,
    kernelVersion: data?.kernelVersion ?? null,
    loading,
    error,
  };
}
