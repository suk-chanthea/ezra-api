# Pagination Feature Documentation

## Overview
This document describes the pagination implementation across all list endpoints in the API. Pagination helps manage large datasets by returning data in smaller, manageable chunks (pages).

## Configuration

- **Default Page Size:** 20 items per page
- **Maximum Page Size:** 100 items per page
- **Minimum Page:** 1

## Pagination Parameters

All `GetAll` endpoints support optional pagination query parameters:

| Parameter | Type | Required | Default | Min | Max | Description |
|-----------|------|----------|---------|-----|-----|-------------|
| `page` | integer | No | 1 | 1 | - | The page number to retrieve |
| `page_size` | integer | No | 20 | 1 | 100 | Number of items per page |

## Response Format

### Without Pagination
When no pagination parameters are provided, the API returns all items in a simple array:

```json
[
  {
    "id": 1,
    "title": "Item 1"
  },
  {
    "id": 2,
    "title": "Item 2"
  }
]
```

### With Pagination
When pagination parameters are provided, the API returns a paginated response with metadata:

```json
{
  "data": [
    {
      "id": 1,
      "title": "Item 1"
    },
    {
      "id": 2,
      "title": "Item 2"
    }
  ],
  "pagination": {
    "current_page": 1,
    "page_size": 20,
    "total_pages": 5,
    "total_records": 95,
    "has_next_page": true,
    "has_prev_page": false
  }
}
```

## Pagination Metadata

| Field | Type | Description |
|-------|------|-------------|
| `current_page` | integer | Current page number |
| `page_size` | integer | Number of items per page |
| `total_pages` | integer | Total number of pages available |
| `total_records` | integer | Total number of records in the dataset |
| `has_next_page` | boolean | Whether there is a next page available |
| `has_prev_page` | boolean | Whether there is a previous page available |

## Supported Endpoints

### 1. Music Endpoints

#### Get All Music (Paginated)
```bash
GET /api/musics?page=1&page_size=20
```

**Example Request:**
```bash
curl -X GET "http://localhost:8080/api/musics?page=1&page_size=20" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Example Response:**
```json
{
  "data": [
    {
      "id": 1,
      "title": "Song Title",
      "cover": "cover.jpg",
      "audio": "audio.mp3",
      "user_id": 5,
      "created_at": "2025-10-24T10:00:00Z",
      "updated_at": "2025-10-24T10:00:00Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "page_size": 20,
    "total_pages": 3,
    "total_records": 45,
    "has_next_page": true,
    "has_prev_page": false
  }
}
```

### 2. Event Endpoints

#### Get All Events (Paginated)
```bash
GET /api/events?page=1&page_size=20
```

**Example Request:**
```bash
curl -X GET "http://localhost:8080/api/events?page=2&page_size=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Example Response:**
```json
{
  "data": [
    {
      "id": 11,
      "title": "Event Title",
      "content": "Event description",
      "cover": "event_cover.jpg",
      "location": "Event Location",
      "start_time": "2025-11-01T18:00:00Z",
      "end_time": "2025-11-01T22:00:00Z",
      "user_id": 3,
      "musics": [],
      "created_at": "2025-10-24T10:00:00Z",
      "updated_at": "2025-10-24T10:00:00Z"
    }
  ],
  "pagination": {
    "current_page": 2,
    "page_size": 10,
    "total_pages": 5,
    "total_records": 50,
    "has_next_page": true,
    "has_prev_page": true
  }
}
```

### 3. Booking Endpoints

#### Get All Bookings (Paginated)
```bash
GET /api/bookings?page=1&page_size=20
```

**Example Request:**
```bash
curl -X GET "http://localhost:8080/api/bookings?page=1&page_size=15" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Example Response:**
```json
{
  "data": [
    {
      "id": 1,
      "event_id": 5,
      "user_id": 3,
      "status": "confirmed",
      "notes": "Looking forward to it!",
      "event": {
        "id": 5,
        "title": "Event Name"
      },
      "user": {
        "id": 3,
        "username": "john_doe"
      },
      "created_at": "2025-10-24T10:00:00Z",
      "updated_at": "2025-10-24T10:00:00Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "page_size": 15,
    "total_pages": 2,
    "total_records": 28,
    "has_next_page": true,
    "has_prev_page": false
  }
}
```

### 4. Favorites Endpoints

#### Get User Favorites (Paginated)
```bash
GET /api/favorites?page=1&page_size=20
```

**Example Request:**
```bash
curl -X GET "http://localhost:8080/api/favorites?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Example Response:**
```json
{
  "data": [
    {
      "id": 1,
      "title": "Favorite Song",
      "cover": "cover.jpg",
      "audio": "audio.mp3",
      "user_id": 2,
      "created_at": "2025-10-24T10:00:00Z",
      "updated_at": "2025-10-24T10:00:00Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "page_size": 10,
    "total_pages": 3,
    "total_records": 25,
    "has_next_page": true,
    "has_prev_page": false
  }
}
```

## Usage Examples

### JavaScript/TypeScript (Fetch API)
```javascript
async function getMusic(page = 1, pageSize = 20) {
  const response = await fetch(
    `http://localhost:8080/api/musics?page=${page}&page_size=${pageSize}`,
    {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    }
  );
  
  const result = await response.json();
  
  console.log(`Page ${result.pagination.current_page} of ${result.pagination.total_pages}`);
  console.log(`Total records: ${result.pagination.total_records}`);
  
  return result.data;
}

// Get first page with 20 items
const musics = await getMusic(1, 20);

// Get second page with 10 items
const moreMusics = await getMusic(2, 10);
```

### Python (Requests)
```python
import requests

def get_events(page=1, page_size=20, token=""):
    url = f"http://localhost:8080/api/events"
    params = {"page": page, "page_size": page_size}
    headers = {"Authorization": f"Bearer {token}"}
    
    response = requests.get(url, params=params, headers=headers)
    result = response.json()
    
    print(f"Page {result['pagination']['current_page']} of {result['pagination']['total_pages']}")
    print(f"Total records: {result['pagination']['total_records']}")
    
    return result['data']

# Get first page
events = get_events(page=1, page_size=20, token="your_token")
```

### cURL
```bash
# Get first page (default 20 items)
curl -X GET "http://localhost:8080/api/musics?page=1" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Get second page with 10 items
curl -X GET "http://localhost:8080/api/musics?page=2&page_size=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Get without pagination (all items)
curl -X GET "http://localhost:8080/api/musics" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Best Practices

### 1. Always Use Pagination for Large Datasets
```bash
# Good: Use pagination for better performance
GET /api/musics?page=1&page_size=20

# Bad: Loading all items at once (can be slow)
GET /api/musics
```

### 2. Handle Pagination in UI
```javascript
function PaginatedList({ endpoint }) {
  const [page, setPage] = useState(1);
  const [data, setData] = useState([]);
  const [pagination, setPagination] = useState(null);

  useEffect(() => {
    fetch(`${endpoint}?page=${page}&page_size=20`)
      .then(res => res.json())
      .then(result => {
        setData(result.data);
        setPagination(result.pagination);
      });
  }, [page]);

  return (
    <div>
      {data.map(item => <Item key={item.id} {...item} />)}
      
      <div className="pagination">
        <button 
          disabled={!pagination?.has_prev_page}
          onClick={() => setPage(page - 1)}
        >
          Previous
        </button>
        
        <span>Page {pagination?.current_page} of {pagination?.total_pages}</span>
        
        <button 
          disabled={!pagination?.has_next_page}
          onClick={() => setPage(page + 1)}
        >
          Next
        </button>
      </div>
    </div>
  );
}
```

### 3. Cache Pagination Results
```javascript
const cache = new Map();

async function getCachedData(page, pageSize) {
  const key = `${page}-${pageSize}`;
  
  if (cache.has(key)) {
    return cache.get(key);
  }
  
  const data = await fetchData(page, pageSize);
  cache.set(key, data);
  
  return data;
}
```

### 4. Validate Page Numbers
```javascript
function validatePage(page, totalPages) {
  if (page < 1) return 1;
  if (page > totalPages) return totalPages;
  return page;
}
```

## Error Handling

### Invalid Pagination Parameters
**Request:**
```bash
GET /api/musics?page=-1&page_size=200
```

**Response (400 Bad Request):**
```json
{
  "error": "invalid pagination parameters"
}
```

The API automatically corrects invalid values:
- Page less than 1 → defaults to 1
- Page size less than 1 → defaults to 20
- Page size greater than 100 → caps at 100

## Performance Considerations

1. **Database Indexing:** All tables use proper indexes for efficient pagination queries
2. **Count Optimization:** Total count is calculated separately for accuracy
3. **Preloading:** Related entities (e.g., musics in events) are preloaded in a single query
4. **Memory Efficiency:** Only requested page is loaded into memory

## Architecture

The pagination feature follows the Clean Architecture pattern:

1. **DTO Layer:** `PaginationRequest` and `PaginationMetadata` DTOs
2. **Repository Layer:** All repositories support `FindAllPaginated(offset, limit)` methods
3. **Use Case Layer:** All use cases provide paginated versions of list methods
4. **Handler Layer:** Handlers parse query parameters and call appropriate use case methods

## Migration Notes

The pagination feature is **backward compatible**:
- Existing endpoints without parameters continue to work
- No database schema changes required
- Old clients can continue using unpaginated endpoints

## Testing

### Manual Testing
```bash
# Test default pagination
curl "http://localhost:8080/api/musics?page=1"

# Test custom page size
curl "http://localhost:8080/api/musics?page=1&page_size=5"

# Test last page
curl "http://localhost:8080/api/musics?page=999&page_size=20"

# Test without pagination (backward compatibility)
curl "http://localhost:8080/api/musics"
```

### Expected Behaviors
- ✅ Returns 20 items by default when pagination is used
- ✅ Returns all items when no pagination parameters provided
- ✅ Caps page_size at 100
- ✅ Defaults to page 1 if page < 1
- ✅ Returns empty array if page exceeds total pages
- ✅ Includes accurate pagination metadata

## Future Enhancements

Potential improvements:
1. Add cursor-based pagination for real-time data
2. Add sorting options (order_by parameter)
3. Add filtering options in combination with pagination
4. Add pagination links (first, last, next, prev URLs)
5. Add ETag/Last-Modified headers for caching
6. Add total count caching for frequently accessed endpoints

