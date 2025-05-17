<script>
  import { shifts, loadShifts, requestShiftAction, workerId } from './lib/stores';
  import { deleteShift } from './services/api';
  
  let currentTab = 'assigned';
  
  // Load shifts on mount
  loadShifts();
</script>

<main>
  <div class="tabs">
    <button class:active={currentTab === 'assigned'} on:click={() => currentTab = 'assigned'}>
      My Shifts ({$shifts.filter(s => s.worker_id === $workerId).length})
    </button>
    <button class:active={currentTab === 'available'} on:click={() => currentTab = 'available'}>
      Available Shifts ({$shifts.filter(s => !s.assigned).length})
    </button>
  </div>

  {#if currentTab === 'assigned'}
    <div class="shift-list">
      {#each $shifts.filter(s => s.assigned) as shift (shift.id)}
        <div class="shift-card {shift.status.toLowerCase()}">
          <div class="header">
            <h3>{shift.date} {shift.start_time} - {shift.end_time}</h3>
            <span class="status">{shift.status || 'approved'}</span>
          </div>
          <p>Role: {shift.role}</p>
          <button class="delete-btn" on:click={() => deleteShift(shift.id)}>Cancel</button>
        </div>
      {/each}
    </div>
  {:else}
    <div class="shift-list">
      {#each $shifts.filter(s => !s.assigned) as shift (shift.id)}
        <div class="shift-card">
          <div class="header">
            <h3>{shift.date} {shift.start_time} - {shift.end_time}</h3>
            <span class="role-badge">{shift.role}</span>
          </div>
          <button 
            class="request-btn" 
            on:click={async () => {
              try {
                await requestShiftAction(shift.id);
              } catch (error) {
                console.error('Shift request failed:', error);
              }
            }}
            disabled={shift.status === 'pending' || shift.assigned}>
            {shift.status === 'pending' ? 'Requested' : 'Request Shift'}
          </button>
        </div>
      {/each}
    </div>
  {/if}
</main>