import {
  QueryClient,
  QueryClientProvider,
  useQuery,
} from "@tanstack/react-query";

interface Product {
  id: number;
  name: string;
  description: string;
  price: number;
  sku: string;
}

function getProducts(): Promise<Product[]> {
  return fetch("http://localhost:9000/").then(
    (res) => res.json() as Promise<Product[]>
  );
}

function Products() {
  const productsQuery = useQuery({
    queryKey: ["products"],
    queryFn: getProducts,
  });

  if (productsQuery.isLoading) {
    return <div>Loading Products...</div>;
  }

  const products = productsQuery.data ?? [];

  return (
    <div>
      {products.map((p) => (
        <div key={p.id}>
          <ul>
            <li>{p.name}</li>
            <li>{p.price}</li>
          </ul>
        </div>
      ))}
    </div>
  );
}

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Products />
    </QueryClientProvider>
  );
}

export default App;
