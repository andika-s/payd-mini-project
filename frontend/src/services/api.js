export const getItems = async () => {
  try {
    const response = await fetch('/api/items');
    return await response.json();
  } catch (error) {
    console.error('Fetch failed:', error);
    return [];
  }
};

export const deleteItem = async (id) => {
  try {
    await fetch(`/api/items/${id}`, { method: 'DELETE' });
  } catch (error) {
    console.error('Delete failed:', error);
  }
};

export const getShifts = async () => {
  try {
    const response = await fetch('http://localhost:8080/api/v1/shifts');
    const { data } = await response.json();
    return data.filter(shift => shift.status !== 'rejected'); // Filter out rejected shifts
  } catch (error) {
    console.error('Fetch failed:', error);
    return [];
  }
};

export const requestShift = async (shiftId, workerId) => {
  try {
    const response = await fetch(`http://localhost:8080/api/v1/shift/request`, {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({ 
        shift_id: shiftId, 
        worker_id: workerId 
      })
    });
    const { data } = await response.json();
    return data;
  } catch (error) {
    console.error('Request failed:', error);
    throw error;
  }
};

export const deleteShift = async (id) => {
  try {
    await fetch(`/api/v1/shift/${id}`, { method: 'DELETE' });
  } catch (error) {
    console.error('Delete failed:', error);
  }
};