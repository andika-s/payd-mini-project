import { writable } from 'svelte/store';

export const shifts = writable([]);
export const workerId = writable(1000); // Mock worker ID

export const loadShifts = async () => {
  const data = await getShifts();
  shifts.set(data.data);
};

export const requestShiftAction = async (shiftId) => {
  const response = await requestShift(shiftId, $workerId);
  if (response && response.data) {
    shifts.update(current => 
      current.map(s => s.id === shiftId ? {...s, ...response.data} : s)
    );
  }
};

// Add initial dummy data for testing
export const items = writable([
  { id: 1, name: 'Test Item', description: 'Sample description' }
]);
export const addItem = (item) => {
  items.update(current => {
    const newItem = item || { name: '', description: '' }; // Replace logical OR assignment
    return [...current, { ...newItem, id: Date.now() }];
  });
};
export const updateItem = (updatedItem) => {
  items.update(items => items.map(item => item.id === updatedItem.id ? updatedItem : item));
};